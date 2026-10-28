'use client';

import { useEffect, useState } from 'react';
import { Users, Shield, UserCheck, UserX, Mail } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Label } from '@/components/ui/label';
import { Spinner } from '@/components/ui/spinner';
import { DataTable, Column } from '@/components/admin/data-table';
import { useAdminStore } from '@/lib/stores/admin-store';
import type { User } from '@/lib/types';
import { toast } from 'sonner';

const columns: Column<User>[] = [
  {
    key: 'name',
    header: 'User',
    sortable: true,
    render: (user) => (
      <div className="flex items-center gap-3">
        <div className="flex h-10 w-10 items-center justify-center rounded-full bg-primary/10">
          {user.profilePic ? (
            <img
              src={user.profilePic}
              alt={`${user.firstName} ${user.lastName}`}
              className="h-10 w-10 rounded-full object-cover"
            />
          ) : (
            <span className="text-sm font-medium text-primary">
              {user.firstName[0]}
              {user.lastName[0]}
            </span>
          )}
        </div>
        <div>
          <p className="font-medium">
            {user.firstName} {user.lastName}
          </p>
          <p className="flex items-center gap-1 text-sm text-muted-foreground">
            <Mail className="h-3 w-3" />
            {user.email}
          </p>
        </div>
      </div>
    ),
  },
  {
    key: 'role',
    header: 'Role',
    sortable: true,
    render: (user) => (
      <span
        className={`rounded-full px-2 py-1 text-xs font-medium ${
          user.role === 'admin'
            ? 'bg-danger-50 text-danger-400'
            : 'bg-primary/10 text-primary'
        }`}
      >
        {user.role}
      </span>
    ),
  },
  {
    key: 'nationality',
    header: 'Nationality',
    render: (user) => user.nationality || '-',
  },
  {
    key: 'isActive',
    header: 'Status',
    render: (user) => (
      <span
        className={`rounded-full px-2 py-1 text-xs font-medium ${
          user.isActive
            ? 'bg-success-50 text-success-400'
            : 'bg-muted text-muted-foreground'
        }`}
      >
        {user.isActive ? 'Active' : 'Inactive'}
      </span>
    ),
  },
];

export default function AdminUsersPage() {
  const { users, isLoading, fetchUsers, updateUserRole, toggleUserStatus } =
    useAdminStore();

  const [isRoleDialogOpen, setIsRoleDialogOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [newRole, setNewRole] = useState('user');

  useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);

  const handleEditRole = (user: User) => {
    setSelectedUser(user);
    setNewRole(user.role);
    setIsRoleDialogOpen(true);
  };

  const handleToggleStatus = async (user: User) => {
    try {
      await toggleUserStatus(user.id);
      toast.success(
        `User ${user.isActive ? 'deactivated' : 'activated'} successfully`
      );
    } catch {
      toast.error('Failed to update user status');
    }
  };

  const handleUpdateRole = async () => {
    if (selectedUser) {
      try {
        await updateUserRole(selectedUser.id, newRole);
        toast.success('User role updated successfully');
        setIsRoleDialogOpen(false);
      } catch {
        toast.error('Failed to update user role');
      }
    }
  };

  const adminCount = users.filter((u) => u.role === 'admin').length;
  const activeCount = users.filter((u) => u.isActive).length;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-foreground">Users</h1>
        <p className="mt-1 text-muted-foreground">
          Manage user accounts and permissions
        </p>
      </div>

      {/* Stats */}
      <div className="grid gap-4 sm:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Total Users
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-2 text-2xl font-bold">
              <Users className="h-5 w-5 text-muted-foreground" />
              {users.length}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Admins
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-2 text-2xl font-bold text-danger-400">
              <Shield className="h-5 w-5" />
              {adminCount}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Active
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-2 text-2xl font-bold text-success-400">
              <UserCheck className="h-5 w-5" />
              {activeCount}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Inactive
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-2 text-2xl font-bold text-muted-foreground">
              <UserX className="h-5 w-5" />
              {users.length - activeCount}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Table */}
      {isLoading && users.length === 0 ? (
        <div className="flex h-64 items-center justify-center">
          <Spinner className="h-8 w-8 text-primary" />
        </div>
      ) : (
        <DataTable
          data={users}
          columns={columns}
          searchPlaceholder="Search users..."
          searchKey="email"
          onEdit={handleEditRole}
          onDelete={handleToggleStatus}
        />
      )}

      {/* Role Dialog */}
      <Dialog open={isRoleDialogOpen} onOpenChange={setIsRoleDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Change User Role</DialogTitle>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <div className="flex items-center gap-3">
              <div className="flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
                <span className="text-lg font-medium text-primary">
                  {selectedUser?.firstName[0]}
                  {selectedUser?.lastName[0]}
                </span>
              </div>
              <div>
                <p className="font-medium">
                  {selectedUser?.firstName} {selectedUser?.lastName}
                </p>
                <p className="text-sm text-muted-foreground">
                  {selectedUser?.email}
                </p>
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="role">Role</Label>
              <Select value={newRole} onValueChange={setNewRole}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="user">User</SelectItem>
                  <SelectItem value="admin">Admin</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsRoleDialogOpen(false)}>
              Cancel
            </Button>
            <Button onClick={handleUpdateRole} disabled={isLoading}>
              {isLoading ? 'Updating...' : 'Update Role'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
