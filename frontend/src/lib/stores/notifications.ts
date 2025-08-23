import { writable } from 'svelte/store';

export interface Notification {
  id: string;
  type: 'success' | 'error' | 'warning' | 'info';
  title: string;
  message?: string;
  duration?: number;
}

const createNotificationStore = () => {
  const { subscribe, update } = writable<Notification[]>([]);

  return {
    subscribe,
    add: (notification: Omit<Notification, 'id'>) => {
      const id = crypto.randomUUID();
      const newNotification = { ...notification, id };

      update((notifications) => [...notifications, newNotification]);

      // Auto remove after duration (default 5 seconds)
      setTimeout(() => {
        update((notifications) => notifications.filter((n) => n.id !== id));
      }, notification.duration || 5000);

      return id;
    },
    remove: (id: string) => {
      update((notifications) => notifications.filter((n) => n.id !== id));
    },
    clear: () => {
      update(() => []);
    }
  };
};

export const notifications = createNotificationStore();
