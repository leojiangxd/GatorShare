import React, { useEffect, useState } from "react";
import NavBar from "./components/NavBar";
import { CheckCircle, AlertCircle } from "lucide-react";
import dayjs from "dayjs";
import { getCsrfToken } from "../utils/functions";
import axios from "axios";

const Notifications = () => {
  const [notifications, setNotifications] = useState([]);
  const apiBaseUrl = import.meta.env.VITE_API_BASE_URL;

  useEffect(() => {
    const fetchNotifications = async () => {
      try {
        const csrfToken = getCsrfToken();
        if (!csrfToken) {
          console.error("CSRF token is missing.");
          navigate("/login");
          return;
        }
        const params = {
          column: "created_at",
          order: "desc",
          limit: 0,
          offset: 0,
          search_key: "",
        };
        const response = await axios.get(`${apiBaseUrl}/api/v1/notification`, {
          params,
          headers: {
            "X-CSRF-Token": csrfToken || "",
          },
          withCredentials: true,
        });
        const notificationsData = response.data.data;
        const sortedNotifications = notificationsData.sort(
          (a, b) => new Date(b.CreatedAt) - new Date(a.CreatedAt)
        );
        setNotifications(sortedNotifications);
      } catch (error) {
        console.error(
          "Error fetching notifications:",
          error.response?.data || error.message
        );
      }
    };

    fetchNotifications();
  }, [apiBaseUrl]);

  const groupNotificationsByDate = (notifications) => {
    return notifications.reduce((groups, notification) => {
      const dateKey = dayjs(notification.CreatedAt).format("YYYY-MM-DD");
      if (!groups[dateKey]) {
        groups[dateKey] = [];
      }
      groups[dateKey].push(notification);
      return groups;
    }, {});
  };

  const groupedNotifications = groupNotificationsByDate(notifications);

  const deleteNotification = async (id) => {
    try {
      const csrfToken = getCsrfToken();
      if (!csrfToken) {
        console.error("CSRF token is missing.");
        return;
      }
      await axios.delete(`${apiBaseUrl}/api/v1/notification/${id}`, {
        headers: {
          "X-CSRF-Token": csrfToken || "",
        },
        withCredentials: true,
      });
      setNotifications((prev) =>
        prev.filter((notification) => notification.id !== id)
      );
    } catch (error) {
      console.error(
        "Error deleting notification:",
        error.response?.data || error.message
      );
    }
  }


  const toggleRead = async (id) => {
    try {
      const csrfToken = getCsrfToken();
      if (!csrfToken) {
        console.error("CSRF token is missing.");
        return;
      }
      const notification = notifications.find((notif) => notif.id === id);
      const updatedReadStatus = !notification.read;

      await axios.put(
        `${apiBaseUrl}/api/v1/notification/${id}`,
        { read: updatedReadStatus },
        {
          headers: {
            "X-CSRF-Token": csrfToken || "",
          },
          withCredentials: true,
        }
      );

      setNotifications((prev) =>
        prev.map((notif) =>
          notif.id === id ? { ...notif, read: updatedReadStatus } : notif
        )
      );
    } catch (error) {
      console.error(
        "Error toggling read status:",
        error.response?.data || error.message
      );
    }
  };

  const markAllAsRead = async () => {
    try {
      const csrfToken = getCsrfToken();
      if (!csrfToken) {
        console.error("CSRF token is missing.");
        return;
      }

      await axios.put(
        `${apiBaseUrl}/api/v1/notification`,
        { read: true },
        {
          headers: {
            "X-CSRF-Token": csrfToken || "",
          },
          withCredentials: true,
        }
      );

      setNotifications((prev) =>
        prev.map((notif) => ({ ...notif, read: true }))
      );
    } catch (error) {
      console.error(
        "Error marking all notifications as read:",
        error.response?.data || error.message
      );
    }
  };

  const markAllAsUnread = async () => {
    try {
      const csrfToken = getCsrfToken();
      if (!csrfToken) {
        console.error("CSRF token is missing.");
        return;
      }

      await axios.put(
        `${apiBaseUrl}/api/v1/notification`,
        { read: true },
        {
          headers: {
            "X-CSRF-Token": csrfToken || "",
          },
          withCredentials: true,
        }
      );

      setNotifications((prev) =>
        prev.map((notif) => ({ ...notif, read: false }))
      );
    } catch (error) {
      console.error(
        "Error marking all notifications as read:",
        error.response?.data || error.message
      );
    }
  };

  return (
    <div className="flex flex-col h-screen bg-base-200">
      <NavBar />
      <div className="p-10 overflow-y-auto flex-grow">
        <div className="flex justify-between items-center mb-6">
          <h1 className="text-3xl font-bold">Notifications</h1>
          <div className="flex gap-2">
            <button
              onClick={markAllAsRead}
              className="btn btn-secondary btn-sm"
            >
              Mark All as Read
            </button>
            <button
              onClick={markAllAsUnread}
              className="btn btn-secondary btn-sm"
            >
              Mark All as Unread
            </button>
          </div>
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
                      <div className="text-sm text-gray-500">
                        {dayjs(notification.CreatedAt).format("h:mm A")}
                      </div>
                      <div className="flex gap-2">
                      <button
                          onClick={() => deleteNotification(notification.id)}
                          className="btn btn-outline btn-xs"
                        >
                          Delete
                        </button>
                        <button
                          onClick={() => toggleRead(notification.id)}
                          className="btn btn-outline btn-xs"
                        >
                          {notification.read
                            ? "Mark as Unread"
                            : "Mark as Read"}
                        </button>
                      </div>
                    </div>
                    <p className="text-lg">{notification.title}</p>
                    <p className="text-sm">{notification.content}</p>
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
