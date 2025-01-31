'use client';

import { Dialog, DialogContent, DialogHeader, DialogFooter, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";

interface ConfirmationDialogProps {
  open: boolean;
  title: string;
  description?: string;
  onConfirm: () => void;
  onCancel: () => void;
}

export default function ConfirmationDialog({ open, title, description, onConfirm, onCancel }: ConfirmationDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onCancel}>
      <DialogContent>
        <DialogTitle>
            {title}
        </DialogTitle>
        <DialogHeader>
          {description && <p className="text-gray-500">{description}</p>}
        </DialogHeader>
        <DialogFooter className="flex justify-end space-x-2">
          <Button variant="outline" onClick={onCancel}>Cancel</Button>
          <Button className="bg-blue-700 hover:bg-blue-900" onClick={onConfirm}>Confirm</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
