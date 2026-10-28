'use client';

import { useEffect, useState } from 'react';
import { useSearchParams } from 'next/navigation';
import { Plus, Package, Star, Calendar } from 'lucide-react';
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
import type { Package as TravelPackage } from '@/lib/types';
import { toast } from 'sonner';
import { format } from 'date-fns';

const columns: Column<TravelPackage>[] = [
  {
    key: 'name',
    header: 'Package',
    sortable: true,
    render: (pkg) => (
      <div className="flex items-center gap-3">
        <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-teal-100">
          <Package className="h-5 w-5 text-teal-600" />
        </div>
        <div>
          <p className="font-medium">{pkg.name}</p>
          <p className="text-sm text-muted-foreground">{pkg.destination}</p>
        </div>
      </div>
    ),
  },
  {
    key: 'durationDays',
    header: 'Duration',
    sortable: true,
    render: (pkg) => (
      <div className="flex items-center gap-1">
        <Calendar className="h-4 w-4 text-muted-foreground" />
        <span>{pkg.durationDays} days</span>
      </div>
    ),
  },
  {
    key: 'basePrice',
    header: 'Price',
    sortable: true,
    render: (pkg) => (
      <span className="font-semibold text-primary">
        {pkg.currency} {pkg.basePrice.toLocaleString()}
      </span>
    ),
  },
  {
    key: 'rating',
    header: 'Rating',
    sortable: true,
    render: (pkg) => (
      <div className="flex items-center gap-1">
        <Star className="h-4 w-4 fill-amber-400 text-amber-400" />
        <span className="font-medium">{pkg.rating.toFixed(1)}</span>
      </div>
    ),
  },
  {
    key: 'maxTravelers',
    header: 'Max Travelers',
    render: (pkg) => <span>{pkg.maxTravelers} people</span>,
  },
];

export default function AdminPackagesPage() {
  const searchParams = useSearchParams();
  const {
    packages,
    isLoading,
    fetchPackages,
    createPackage,
    updatePackage,
    deletePackage,
  } = useAdminStore();

  const [isFormOpen, setIsFormOpen] = useState(false);
  const [isDeleteOpen, setIsDeleteOpen] = useState(false);
  const [selectedItem, setSelectedItem] = useState<TravelPackage | null>(null);
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    destination: '',
    durationDays: 7,
    startDateFrom: '',
    startDateTo: '',
    maxTravelers: 20,
    basePrice: 0,
    currency: 'USD',
    tags: '',
  });

  useEffect(() => {
    fetchPackages();
  }, [fetchPackages]);

  useEffect(() => {
    if (searchParams.get('action') === 'new') {
      setIsFormOpen(true);
    }
  }, [searchParams]);

  const handleCreate = () => {
    setSelectedItem(null);
    setFormData({
      name: '',
      description: '',
      destination: '',
      durationDays: 7,
      startDateFrom: '',
      startDateTo: '',
      maxTravelers: 20,
      basePrice: 0,
      currency: 'USD',
      tags: '',
    });
    setIsFormOpen(true);
  };

  const handleEdit = (item: TravelPackage) => {
    setSelectedItem(item);
    setFormData({
      name: item.name,
      description: item.description,
      destination: item.destination,
      durationDays: item.durationDays,
      startDateFrom: item.startDateFrom?.slice(0, 10) || '',
      startDateTo: item.startDateTo?.slice(0, 10) || '',
      maxTravelers: item.maxTravelers,
      basePrice: item.basePrice,
      currency: item.currency,
      tags: item.tags?.join(', ') || '',
    });
    setIsFormOpen(true);
  };

  const handleDelete = (item: TravelPackage) => {
    setSelectedItem(item);
    setIsDeleteOpen(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const data: Partial<TravelPackage> = {
      name: formData.name,
      description: formData.description,
      destination: formData.destination,
      durationDays: formData.durationDays,
      startDateFrom: formData.startDateFrom,
      startDateTo: formData.startDateTo,
      maxTravelers: formData.maxTravelers,
      basePrice: formData.basePrice,
      currency: formData.currency,
      tags: formData.tags.split(',').map((t) => t.trim()).filter(Boolean),
      includedComponents: [],
      images: [],
      rating: 0,
    };

    try {
      if (selectedItem) {
        await updatePackage(selectedItem.id, data);
        toast.success('Package updated successfully');
      } else {
        await createPackage(data as Omit<TravelPackage, 'id'>);
        toast.success('Package created successfully');
      }
      setIsFormOpen(false);
    } catch {
      toast.error('Failed to save package');
    }
  };

  const handleConfirmDelete = async () => {
    if (selectedItem) {
      try {
        await deletePackage(selectedItem.id);
        toast.success('Package deleted successfully');
        setIsDeleteOpen(false);
      } catch {
        toast.error('Failed to delete package');
      }
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Packages</h1>
          <p className="mt-1 text-muted-foreground">
            Manage travel packages and bundles
          </p>
        </div>
        <Button onClick={handleCreate} className="gap-2">
          <Plus className="h-4 w-4" />
          Add Package
        </Button>
      </div>

      {/* Stats */}
      <div className="grid gap-4 sm:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Total Packages
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{packages.length}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Avg Duration
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {packages.length > 0
                ? Math.round(
                    packages.reduce((sum, p) => sum + p.durationDays, 0) /
                      packages.length
                  )
                : 0}{' '}
              days
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Avg Price
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-primary">
              $
              {packages.length > 0
                ? Math.round(
                    packages.reduce((sum, p) => sum + p.basePrice, 0) /
                      packages.length
                  ).toLocaleString()
                : 0}
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
              {packages.length > 0
                ? (
                    packages.reduce((sum, p) => sum + p.rating, 0) /
                    packages.length
                  ).toFixed(1)
                : '0.0'}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Table */}
      {isLoading && packages.length === 0 ? (
        <div className="flex h-64 items-center justify-center">
          <Spinner className="h-8 w-8 text-primary" />
        </div>
      ) : (
        <DataTable
          data={packages}
          columns={columns}
          searchPlaceholder="Search packages..."
          searchKey="name"
          onEdit={handleEdit}
          onDelete={handleDelete}
        />
      )}

      {/* Form Dialog */}
      <Dialog open={isFormOpen} onOpenChange={setIsFormOpen}>
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <Package className="h-5 w-5" />
              {selectedItem ? 'Edit Package' : 'Create Package'}
            </DialogTitle>
          </DialogHeader>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">Package Name</Label>
              <Input
                id="name"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                placeholder="European Adventure"
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
                placeholder="Describe the package..."
                rows={3}
              />
            </div>
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="destination">Destination</Label>
                <Input
                  id="destination"
                  value={formData.destination}
                  onChange={(e) =>
                    setFormData({ ...formData, destination: e.target.value })
                  }
                  placeholder="Paris, France"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="durationDays">Duration (days)</Label>
                <Input
                  id="durationDays"
                  type="number"
                  min="1"
                  value={formData.durationDays}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      durationDays: parseInt(e.target.value),
                    })
                  }
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="startDateFrom">Available From</Label>
                <Input
                  id="startDateFrom"
                  type="date"
                  value={formData.startDateFrom}
                  onChange={(e) =>
                    setFormData({ ...formData, startDateFrom: e.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="startDateTo">Available To</Label>
                <Input
                  id="startDateTo"
                  type="date"
                  value={formData.startDateTo}
                  onChange={(e) =>
                    setFormData({ ...formData, startDateTo: e.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="maxTravelers">Max Travelers</Label>
                <Input
                  id="maxTravelers"
                  type="number"
                  min="1"
                  value={formData.maxTravelers}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      maxTravelers: parseInt(e.target.value),
                    })
                  }
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="basePrice">Base Price</Label>
                <Input
                  id="basePrice"
                  type="number"
                  min="0"
                  step="0.01"
                  value={formData.basePrice}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      basePrice: parseFloat(e.target.value),
                    })
                  }
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="currency">Currency</Label>
                <Input
                  id="currency"
                  value={formData.currency}
                  onChange={(e) =>
                    setFormData({ ...formData, currency: e.target.value })
                  }
                  placeholder="USD"
                />
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="tags">Tags (comma-separated)</Label>
              <Input
                id="tags"
                value={formData.tags}
                onChange={(e) => setFormData({ ...formData, tags: e.target.value })}
                placeholder="Adventure, Cultural, Family-friendly"
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
        title="Delete Package"
        description={`Are you sure you want to delete "${selectedItem?.name}"? This action cannot be undone.`}
        isLoading={isLoading}
      />
    </div>
  );
}
