import React, { useEffect, useState } from "react";
import NavBar from "./components/NavBar";
import { CheckCircle, AlertCircle } from "lucide-react";
import dayjs from "dayjs";

// Sample notifications data
const sampleNotifications = [
  {
    id: 1,
    message: "Your request for tutoring in Calculus has been accepted.",
    type: "help",
    timestamp: "2025-03-31T14:30:00",
    read: false,
  },
  {
    id: 2,
    message: "New item posted: Gatorshare textbooks for sale.",
    type: "sale",
    timestamp: "2025-03-31T10:15:00",
    read: true,
  },
  {
    id: 3,
    message: "Your post 'Looking for study buddies' got a new comment.",
    type: "activity",
    timestamp: "2025-03-30T16:45:00",
    read: false,
  },
  {
    id: 4,
    message: "Reminder: Upcoming UFL committee meeting tomorrow.",
    type: "reminder",
    timestamp: "2025-03-30T09:00:00",
    read: true,
  },
  {
    id: 5,
    message: "Your sale item 'Old Laptop' has a new offer.",
    type: "sale",
    timestamp: "2025-03-29T12:00:00",
    read: false,
  },
];

const Notifications = () => {
  const [notifications, setNotifications] = useState([]);

  // Simulate fetching notifications from backend
  useEffect(() => {
    // Sort notifications by timestamp (descending)
    const sortedNotifications = sampleNotifications.sort(
      (a, b) => new Date(b.timestamp) - new Date(a.timestamp)
    );
    setNotifications(sortedNotifications);
  }, []);

  // Group notifications by date
  const groupNotificationsByDate = (notifications) => {
    return notifications.reduce((groups, notification) => {
      const dateKey = dayjs(notification.timestamp).format("YYYY-MM-DD");
      if (!groups[dateKey]) {
        groups[dateKey] = [];
      }
      groups[dateKey].push(notification);
      return groups;
    }, {});
  };

  const groupedNotifications = groupNotificationsByDate(notifications);

  const toggleRead = (id) => {
    setNotifications((prev) =>
      prev.map((notif) =>
        notif.id === id ? { ...notif, read: !notif.read } : notif
      )
    );
  };

  const markAllAsRead = () => {
    setNotifications((prev) =>
      prev.map((notif) => ({ ...notif, read: true }))
    );
  };

  return (
    <div className="flex flex-col h-screen bg-base-200">
      <NavBar />
      <div className="p-10 overflow-y-auto flex-grow">
        <div className="flex justify-between items-center mb-6">
          <h1 className="text-3xl font-bold">Notifications</h1>
          <button onClick={markAllAsRead} className="btn btn-secondary btn-sm">
            Mark All as Read
          </button>
        </div>

        {Object.keys(groupedNotifications)
          .sort((a, b) => new Date(b) - new Date(a))
          .map((dateKey) => (
            <div key={dateKey} className="mb-8">
              <h2 className="text-xl font-semibold mb-4">
                {dayjs(dateKey).format("MMMM D, YYYY")}
              </h2>
              <div className="flex flex-col gap-4">
                {groupedNotifications[dateKey].map((notification) => (
                  <div
                    key={notification.id}
                    className={`card ${
                      notification.read
                        ? "bg-base-100"
                        : "bg-base-200 border-l-4 border-primary"
                    } shadow-md p-4`}
                  >
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-3">
                        {notification.type === "sale" ? (
                          <AlertCircle className="w-6 h-6 text-warning" />
                        ) : (
                          <CheckCircle className="w-6 h-6 text-success" />
                        )}
                        <p className="text-lg">{notification.message}</p>
                      </div>
                      <div className="text-sm text-gray-500">
                        {dayjs(notification.timestamp).format("h:mm A")}
                      </div>
                    </div>
                    <div className="mt-2 flex justify-end">
                      <button
                        onClick={() => toggleRead(notification.id)}
                        className="btn btn-outline btn-xs"
                      >
                        {notification.read ? "Mark as Unread" : "Mark as Read"}
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          ))}
      </div>
    </div>
  );
};

export default Notifications;
