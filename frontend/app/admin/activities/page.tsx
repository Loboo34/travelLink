'use client';

import { useEffect, useState } from 'react';
import { useSearchParams } from 'next/navigation';
import { Plus, Ticket, Star, Clock } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Spinner } from '@/components/ui/spinner';
import { DataTable, Column } from '@/components/admin/data-table';
import { DeleteDialog } from '@/components/admin/delete-dialog';
import { useAdminStore } from '@/lib/stores/admin-store';
import type { Activity } from '@/lib/types';
import { toast } from 'sonner';

const columns: Column<Activity>[] = [
  {
    key: 'title',
    header: 'Activity',
    sortable: true,
    render: (activity) => (
      <div className="flex items-center gap-3">
        <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-amber-100">
          <Ticket className="h-5 w-5 text-amber-600" />
        </div>
        <div>
          <p className="font-medium">{activity.title}</p>
          <p className="text-sm text-muted-foreground">
            {activity.categories?.join(', ')}
          </p>
        </div>
      </div>
    ),
  },
  {
    key: 'location',
    header: 'Location',
    render: (activity) => (
      <span>
        {activity.city}, {activity.country}
      </span>
    ),
  },
  {
    key: 'durationMinutes',
    header: 'Duration',
    sortable: true,
    render: (activity) => (
      <div className="flex items-center gap-1">
        <Clock className="h-4 w-4 text-muted-foreground" />
        <span>
          {activity.durationMinutes >= 60
            ? `${Math.floor(activity.durationMinutes / 60)}h ${activity.durationMinutes % 60}m`
            : `${activity.durationMinutes}m`}
        </span>
      </div>
    ),
  },
  {
    key: 'rating',
    header: 'Rating',
    sortable: true,
    render: (activity) => {
      const rating = typeof activity.rating === 'number' ? activity.rating : 0;
      const reviewCount =
        typeof activity.reviewCount === 'number' ? activity.reviewCount : 0;

      return (
        <div className="flex items-center gap-1">
          <Star className="h-4 w-4 fill-amber-400 text-amber-400" />
          <span className="font-medium">{rating.toFixed(1)}</span>
          <span className="text-muted-foreground">({reviewCount})</span>
        </div>
      );
    },
  },
];

export default function AdminActivitiesPage() {
  const searchParams = useSearchParams();
  const {
    activities,
    isLoading,
    fetchActivities,
    createActivity,
    updateActivity,
    deleteActivity,
  } = useAdminStore();

  const [isFormOpen, setIsFormOpen] = useState(false);
  const [isDeleteOpen, setIsDeleteOpen] = useState(false);
  const [selectedItem, setSelectedItem] = useState<Activity | null>(null);
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    city: '',
    country: '',
    meetingPoint: '',
    categories: '',
    durationMinutes: 60,
    inclusions: '',
    exclusions: '',
  });

  useEffect(() => {
    fetchActivities();
  }, [fetchActivities]);

  useEffect(() => {
    if (searchParams.get('action') === 'new') {
      setIsFormOpen(true);
    }
  }, [searchParams]);

  const handleCreate = () => {
    setSelectedItem(null);
    setFormData({
      title: '',
      description: '',
      city: '',
      country: '',
      meetingPoint: '',
      categories: '',
      durationMinutes: 60,
      inclusions: '',
      exclusions: '',
    });
    setIsFormOpen(true);
  };

  const handleEdit = (item: Activity) => {
    setSelectedItem(item);
    setFormData({
      title: item.title,
      description: item.description,
      city: item.city,
      country: item.country,
      meetingPoint: item.meetingPoint,
      categories: item.categories?.join(', ') || '',
      durationMinutes: item.durationMinutes,
      inclusions: item.inclusions?.join(', ') || '',
      exclusions: item.exclusions?.join(', ') || '',
    });
    setIsFormOpen(true);
  };

  const handleDelete = (item: Activity) => {
    setSelectedItem(item);
    setIsDeleteOpen(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const data: Partial<Activity> = {
      title: formData.title,
      description: formData.description,
      city: formData.city,
      country: formData.country,
      location: { latitude: 0, longitude: 0 },
      meetingPoint: formData.meetingPoint,
      categories: formData.categories.split(',').map((c) => c.trim()).filter(Boolean),
      durationMinutes: formData.durationMinutes,
      inclusions: formData.inclusions.split(',').map((i) => i.trim()).filter(Boolean),
      exclusions: formData.exclusions.split(',').map((e) => e.trim()).filter(Boolean),
      images: [],
      rating: 0,
      reviewCount: 0,
    };

    try {
      if (selectedItem) {
        await updateActivity(selectedItem.id, data);
        toast.success('Activity updated successfully');
      } else {
        await createActivity(data as Omit<Activity, 'id'>);
        toast.success('Activity created successfully');
      }
      setIsFormOpen(false);
    } catch {
      toast.error('Failed to save activity');
    }
  };

  const handleConfirmDelete = async () => {
    if (selectedItem) {
      try {
        await deleteActivity(selectedItem.id);
        toast.success('Activity deleted successfully');
        setIsDeleteOpen(false);
      } catch {
        toast.error('Failed to delete activity');
      }
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Activities</h1>
          <p className="mt-1 text-muted-foreground">
            Manage tours, experiences, and activities
          </p>
        </div>
        <Button onClick={handleCreate} className="gap-2">
          <Plus className="h-4 w-4" />
          Add Activity
        </Button>
      </div>

      {/* Stats */}
      <div className="grid gap-4 sm:grid-cols-3">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Total Activities
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{activities.length}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Avg Duration
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-1 text-2xl font-bold">
              <Clock className="h-5 w-5 text-muted-foreground" />
              {activities.length > 0
                ? Math.round(
                    activities.reduce((sum, a) => sum + a.durationMinutes, 0) /
                      activities.length
                  )
                : 0}{' '}
              min
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Avg Rating
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-1 text-2xl font-bold">
              <Star className="h-5 w-5 fill-amber-400 text-amber-400" />
              {activities.length > 0
                ? (
                    activities.reduce(
                      (sum, a) =>
                        sum +
                        (typeof a.rating === 'number' ? a.rating : 0),
                      0
                    ) / activities.length
                  ).toFixed(1)
                : '0.0'}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Table */}
      {isLoading && activities.length === 0 ? (
        <div className="flex h-64 items-center justify-center">
          <Spinner className="h-8 w-8 text-primary" />
        </div>
      ) : (
        <DataTable
          data={activities}
          columns={columns}
          searchPlaceholder="Search activities..."
          searchKey="title"
          onEdit={handleEdit}
          onDelete={handleDelete}
        />
      )}

      {/* Form Dialog */}
      <Dialog open={isFormOpen} onOpenChange={setIsFormOpen}>
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <Ticket className="h-5 w-5" />
              {selectedItem ? 'Edit Activity' : 'Create Activity'}
            </DialogTitle>
          </DialogHeader>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="title">Activity Title</Label>
              <Input
                id="title"
                value={formData.title}
                onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                placeholder="City Walking Tour"
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="description">Description</Label>
              <Textarea
                id="description"
                value={formData.description}
                onChange={(e) =>
                  setFormData({ ...formData, description: e.target.value })
                }
                placeholder="Describe the activity..."
                rows={3}
              />
            </div>
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="city">City</Label>
                <Input
                  id="city"
                  value={formData.city}
                  onChange={(e) => setFormData({ ...formData, city: e.target.value })}
                  placeholder="Paris"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="country">Country</Label>
                <Input
                  id="country"
                  value={formData.country}
                  onChange={(e) =>
                    setFormData({ ...formData, country: e.target.value })
                  }
                  placeholder="France"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="meetingPoint">Meeting Point</Label>
                <Input
                  id="meetingPoint"
                  value={formData.meetingPoint}
                  onChange={(e) =>
                    setFormData({ ...formData, meetingPoint: e.target.value })
                  }
                  placeholder="Eiffel Tower entrance"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="durationMinutes">Duration (minutes)</Label>
                <Input
                  id="durationMinutes"
                  type="number"
                  min="15"
                  value={formData.durationMinutes}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      durationMinutes: parseInt(e.target.value),
                    })
                  }
                  required
                />
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="categories">Categories (comma-separated)</Label>
              <Input
                id="categories"
                value={formData.categories}
                onChange={(e) =>
                  setFormData({ ...formData, categories: e.target.value })
                }
                placeholder="Tours, Walking, Historical"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="inclusions">What&apos;s Included (comma-separated)</Label>
              <Input
                id="inclusions"
                value={formData.inclusions}
                onChange={(e) =>
                  setFormData({ ...formData, inclusions: e.target.value })
                }
                placeholder="Guide, Entrance fees, Lunch"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="exclusions">What&apos;s Not Included (comma-separated)</Label>
              <Input
                id="exclusions"
                value={formData.exclusions}
                onChange={(e) =>
                  setFormData({ ...formData, exclusions: e.target.value })
                }
                placeholder="Transport, Tips, Personal expenses"
              />
            </div>
            <div className="flex justify-end gap-2">
              <Button
                type="button"
                variant="outline"
                onClick={() => setIsFormOpen(false)}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={isLoading}>
                {isLoading ? 'Saving...' : selectedItem ? 'Update' : 'Create'}
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>

      {/* Delete Dialog */}
      <DeleteDialog
        open={isDeleteOpen}
        onOpenChange={setIsDeleteOpen}
        onConfirm={handleConfirmDelete}
        title="Delete Activity"
        description={`Are you sure you want to delete "${selectedItem?.title}"? This action cannot be undone.`}
        isLoading={isLoading}
      />
    </div>
  );
}
