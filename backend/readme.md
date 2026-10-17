# TravelLink Backend

TravelLink is a Go HTTP API for travel discovery and bookings. It provides authentication, flight/accommodation/activity/package search and management, booking and cancellation workflows, reviews, and admin operations. MongoDB is the persistent store.

## Contents

- [Stack and architecture](#stack-and-architecture)
- [Requirements](#requirements)
- [Configuration](#configuration)
- [Run locally](#run-locally)
- [API conventions](#api-conventions)
- [API reference](#api-reference)
- [MongoDB collections and indexes](#mongodb-collections-and-indexes)
- [Seed data](#seed-data)
- [Development and verification](#development-and-verification)
- [Known implementation notes](#known-implementation-notes)

## Stack and architecture

- Go module: `github.com/Loboo34/travel`; `go.mod` declares Go `1.25.5`.
- HTTP routing: `net/http` and Gorilla Mux. Routes are mounted at the root; there is no `/api` prefix.
- Persistence: MongoDB using the official Go driver. The database name is hard-coded to `travellink`.
- Authentication: JWT signed with HS256; user passwords are hashed with bcrypt.
- Images: Cloudinary is used for profile image uploads.
- Logging: Zap, with production formatting enabled when `ENV=production`.
- Payments: the current `NewStripeProvider` returns a local stub provider. Booking charges and refunds are not connected to Stripe or another live payment processor.

Requests flow through handlers, services, and repositories. `models/` contains the MongoDB/JSON data types; `repository/` contains persistence and search operations; `service/` implements domain workflows; `auth/` and `middleware/` provide authentication, authorization, and CORS support.

```text
main.go -> handlers -> service -> repository -> MongoDB
				   -> auth / middleware / utils
```

## Requirements

- Go `1.25.5` or a compatible newer toolchain.
- A reachable MongoDB deployment.
- Cloudinary credentials. The server initializes Cloudinary during startup, so valid values are required even if you do not plan to upload images.

## Configuration

Create `backend/.env` and run the server with `backend/` as the working directory. The server calls `godotenv.Load()` during initialization and exits if it cannot load `.env`.

```env
# MongoDB connection string; database name is always "travellink".
MONGO_URI=mongodb://127.0.0.1:27017

# HTTP listener (defaults to 8080 if omitted).
PORT=8080
ENV=development

# JWT signing secret. Use a long, random value outside local development.
JWT_SECRET=replace-with-a-long-random-secret

# Cloudinary credentials (required at startup).
CLOUDINARY_CLOUD_NAME=your-cloud-name
CLOUDINARY_API_KEY=your-api-key
CLOUDINARY_API_SECRET=your-api-secret

# Bootstrap admin account. If it already exists, startup leaves it unchanged.
FIRSTNAME=Travel
LASTNAME=Admin
EMAIL=admin@example.com
PASSWORD=replace-with-a-strong-password

# Optional comma-separated browser origins. Defaults to localhost:3000 and :5173.
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

`JWT_EXPIRY`, `MONGODB_URI`, and `DB_NAME` are not read by the current server. JWT lifetime is fixed at 720 hours (30 days). Keep `.env` out of version control and do not use development secrets in deployed environments.

## Run locally

From the backend directory:

```sh
go mod download
go run .
```

By default the API listens on `http://localhost:8080`; set `PORT` to change the port. Startup connects to MongoDB, ensures required indexes exist, checks/creates the configured admin user, and initializes Cloudinary before serving requests.

Build the server:

```sh
go build -o travellink-backend .
```

## API conventions

- All paths below are relative to the server root, for example `http://localhost:8080/auth/login`.
- JSON success responses use `{ "success": true, "data": ... }`.
- JSON errors use `{ "success": false, "error": { "coode": ..., "status": ..., "message": ... } }` (`coode` is the current JSON field spelling).
- Protected routes require `Authorization: Bearer <token>`. Registration and login return a token in `data.token`.
- Admin routes require a token whose role is `admin`. A bootstrap admin is created from the admin environment variables above if that email is not already present.
- Search endpoints use query parameters and return JSON. Dates are generally parsed as `YYYY-MM-DD`; see the model and parser for endpoint-specific details.
- CORS allows `http://localhost:3000` and `http://localhost:5173` by default. Override this with comma-separated `ALLOWED_ORIGINS` values.

### Authentication and profile

| Method | Path | Access | Notes |
| --- | --- | --- | --- |
| `POST` | `/auth/register` | Public | JSON: `firstName`, `lastName`, `email`, `password` (minimum 8 characters). |
| `POST` | `/auth/login` | Public | JSON: `email`, `password`. Returns a JWT and user object. |
| `GET` | `/auth/profile` | User | Current user's profile. |
| `PATCH` | `/user/profile` | User | Updates profile fields. |
| `PATCH` | `/user/password` | User | JSON: `currentPassword`, `newPassword`. |
| `PATCH` | `/auth/profilepic` | User | Multipart form field `image`; one image, maximum 10 MB. |

### Public catalog and search

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/flights` | List flights. |
| `GET` | `/flight/{flightID}` | Get a flight. |
| `GET` | `/offers` | List flight offers. |
| `GET` | `/offer/{offerID}` | Get an offer. |
| `GET` | `/accommodations` | List accommodations. |
| `GET` | `/accommodation/{accommodationID}` | Get an accommodation. |
| `GET` | `/available` | List accommodation availability. |
| `GET` | `/available/{availilableID}` | Get an availability record. |
| `GET` | `/activitties` | List activities. The misspelling is part of the registered path. |
| `GET` | `/activity/{activityID}` | Get an activity. |
| `GET` | `/timeslots` | List activity time slots. |
| `GET` | `/timeslot/{timeslotID}` | Get a time slot. |
| `GET` | `/packages` | List packages. |
| `GET` | `/package/{packageID}` | Get a package. |
| `GET` | `/flight/search` | Search flights. Query: `originCode`, `destinationCode`, `departureTime`, `adults`; optional `returnDate`, `children`, `infants`, `cabinClass`, `sortBy`, `page`, `pageSize`. |
| `GET` | `/accommodation/search` | Search stays. Query: `city` or `latitude`/`longitude`, `checkInDate`, `checkOutDate`, `adults`; optional `country`, `radiusKm`, `children`, `infants`, property type, sort and pagination values. |
| `GET` | `/activity/search` | Query: `date`, `adults`; optional `city`, `country`, `latitude`, `longitude`, `radiusKm`, `children`, `infants`, `maxDurationMinutes`, `category`, `sortBy`, `page`, `pageSize`. |
| `GET` | `/package/search` | Query: `destination`, `startDate`, `endDate`, `adults`; optional `children`, `infants`, `maxBudget`, repeated `tag`, `components`, `sortBy`, `page`, `pageSize`. |

Search validation and supported enum values are defined in `models/search.go` and the parsers in `handlers/User/`. Query names are case-sensitive. The accommodation and activity models also support location coordinates and radius searches.

### Booking and reviews

| Method | Path | Access | Purpose |
| --- | --- | --- | --- |
| `POST` | `/flight/book` | User | Create a flight booking. |
| `POST` | `/flight/cancel` | User | Request flight booking cancellation. |
| `POST` | `/accommodation/book` | User | Create an accommodation booking. |
| `POST` | `/accommodation/cancel` | User | Request accommodation booking cancellation. |
| `POST` | `/activity/book` | User | Create an activity booking. |
| `POST` | `/activity/cancel` | User | Request activity booking cancellation. |
| `POST` | `/package/book` | User | Create a package booking. |
| `POST` | `/package/cancel` | User | Request package booking cancellation. |
| `POST` | `/review/create` | User | Create a review. |
| `GET` | `/review/get?referenceID={id}&reviewFor={type}` | User | Fetch reviews for a target. Both query parameters are required. |
| `PATCH` | `/review/update` | User | Update a review; see implementation notes below. |
| `DELETE` | `/review/delete/{reviewID}` | User | Delete a review; see implementation notes below. |

Booking JSON request types are `FlightRequest`, `AccommodationBookingRequest`, `ActivityBookingRequest`, and `PackageBookingRequest` in `models/booking_params.go`. Cancellation requests use `models/cancel_params.go`. Payment method values currently accepted by booking validation are defined in `models/Booking.go`; payment processing itself is stubbed.

### Admin operations

All routes in this section require an admin JWT.

| Method | Path | Purpose |
| --- | --- | --- |
| `POST` | `/flight/add` | Create flight. |
| `PUT` | `/flight/update/{flightID}` | Update flight. |
| `PATCH` | `/flight/status/{flightID}` | Change flight status. |
| `DELETE` | `/flight/delete/{flightID}` | Delete flight. |
| `POST` | `/offer/create` | Create flight offer. |
| `PUT` | `/offer/update/{offerID}` | Update offer. |
| `PATCH` | `/offer/isactive/{offerID}` | Change offer active state. |
| `DELETE` | `/offer/delete/{offerID}` | Delete offer. |
| `POST` | `/accommodation/add` | Create accommodation. |
| `PATCH` | `/accommodation/update/{accommodationID}` | Update accommodation. |
| `DELETE` | `/accommodation/delete/{accommodationID}` | Delete accommodation. |
| `POST` | `/availability/add` | Add accommodation availability. |
| `PATCH` | `/availability/status/{availabilityID}` | Change availability active state. |
| `DELETE` | `/availability/remove/{availabilityID}` | Remove availability. |
| `POST` | `/activity/add` | Create activity. |
| `PUT` or `PATCH` | `/activity/update/{activityID}` | Update activity. |
| `DELETE` | `/activity/delete/{activityID}` | Delete activity. |
| `POST` | `/activity/{activityID}/timeslot` | Add an activity time slot. |
| `POST` | `/package/add` | Create package. |
| `PUT` or `PATCH` | `/package/update/{packageID}` | Update package. |
| `PATCH` | `/package/status/{packageID}` | Change package active state. |
| `DELETE` | `/package/delete/{packageID}` | Delete package. |
| `GET` | `/admin/stats` | Dashboard statistics. |
| `GET` | `/admin/users` | List users. |
| `PUT` | `/admin/user/{userId}/role` | Change a user's role. |
| `PUT` | `/admin/user/{userId}/toggle-status` | Toggle a user's active state. |
| `GET` | `/admin/bookings` | List bookings. |
| `PUT` | `/admin/booking/{bookingId}/status` | Change booking status. |
| `GET` | `/admin/airlines` | List airlines. |
| `GET` | `/admin/airline/{airlineID}` | Get an airline. |
| `GET` | `/admin/airports` | List airports. |
| `GET` | `/admin/airport/{airportID}` | Get an airport. |
| `GET` | `/admin/planes` | List planes. |
| `GET` | `/admin/plane/{planeID}` | Get a plane. |

Although airline, airport, route and plane mutation handlers exist in `handlers/Admin/Airline.go`, only the airline/airport/plane read endpoints listed above are registered in `main.go`; route handlers are not registered.

## MongoDB collections and indexes

The database is `travellink`. Collections are created as needed by MongoDB and include users, flights/offers, airlines, airports, planes, routes, accommodations and availability, activities and time slots, packages, four booking collections, and reviews. Collection names and fields are defined by repositories and model BSON tags.

Startup creates indexes for flight offers, geospatial and location-based accommodation/activity queries, availability/time slots, booking lookups by user and status, unique user email, and review targets/unique user-target pairs. Startup fails if MongoDB is unreachable or index creation fails. Accommodation and activity geo queries require MongoDB geospatial index support.

## Seed data

The optional seed program is `cmd/seed`. Do not run it against a database containing data you need: the current implementation drops seed collections unconditionally in several seed functions, and the `-fresh` flag does not make those operations safe. It also contains a hard-coded MongoDB URI fallback; configure `MONGO_URI` explicitly and review the seed implementation before running it.

```sh
MONGO_URI=mongodb://127.0.0.1:27017 go run ./cmd/seed
```

In PowerShell, set `$env:MONGO_URI` before running the command. Seed data includes airports, airlines, accommodations/availability, and activities/time slots. It does not create initial flights or packages.

## Development and verification

From `backend/`:

```sh
go test ./...
go build ./...
```

Keep HTTP-specific parsing and status codes in handlers, business rules in services, database operations in repositories, and shared request/response types in `models/`. Add or update tests when changing contracts, booking workflows, authorization, or persistence behavior. There are currently no backend `_test.go` files, so `go test ./...` primarily verifies that packages compile.

## Known implementation notes

- The API does not register a health-check endpoint or an `/api` route prefix.
- `GET /activitties` is misspelled in the router; clients must use that exact path unless the route is corrected.
- The flight-search parser currently attempts to parse a supplied `returnDate` with an empty time layout; return-date searches should be treated as unsupported until that parser is fixed.
- Package search reads repeated `tag` values and the `components` parameter; the parser also reads `component` values but then overwrites them, so use `components` for now.
- Review update/delete handlers read a standard-library path value named `id`, while the router uses Gorilla Mux and the delete route names its parameter `reviewID`; update is also registered without an ID in its path. These routes need code changes before they can reliably resolve a review ID.
- Auth responses serialize the `User` model directly, and its JSON representation includes the stored `password` field (bcrypt hash). Avoid exposing these responses to untrusted clients until a safe response DTO omits it; never return or log passwords.
- The seed command's `-fresh` flag is not a reliable safeguard; some collections are dropped regardless of the flag.
- `payment.NewStripeProvider()` is currently a stub: it reports a synthetic paid result and refunds are no-ops. Do not treat it as production payment processing.