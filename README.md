# Zulassungsstelle Bot

A simple, configurable bot for automatically finding and booking appointments at the local registration office (Zulassungsstelle Pinneberg).

This bot automates the tedious process of repeatedly checking the website for an available appointment slot that matches your criteria. It runs in the terminal, asks for your desired dates and times, and then works in the background to secure a spot for you.

## 🎥 Demo

![](demo.gif)

---

## ✨ Features

- **Interactive Terminal UI:** A user-friendly command-line interface to configure your appointment preferences.
- **Automatic Watching:** The bot polls the website at a configurable interval.
- **Flexible Scheduling:** Define specific days, recurring weekdays, and time ranges for your desired appointment.
- **Automatic Booking:** Once a matching slot is found, the bot automatically navigates the booking process, fills in your details, and confirms the appointment.
- **Debug Mode:** Run the bot with a visible browser window to see exactly what it's doing.

---

## 📋 Requirements

- [Go](https://go.dev/doc/install) (version 1.18 or newer)
- A compatible web browser installed (e.g., Google Chrome, Chromium, or Microsoft Edge).

---

## 🚀 How to Run

1.  **Clone the Repository:**

    ```bash
    git clone https://github.com/mlentzler/ZulassungsstelleBot.git
    cd ZulassungsstelleBot
    ```

2.  **Run the Application:**
    Execute the bot from your terminal using the following command:

    ```bash
    go run ./cmd/zulassungsstellebot
    ```

3.  **Follow the TUI:**
    The bot will start an interactive setup process in your terminal. Follow the prompts to select the correct office, service, and your desired appointment times. You will also be asked to enter your personal details (Name, Email, Phone) required for the booking.

4.  **Let it Run:**
    Once configured, the bot will start watching the website. You can leave the terminal window running in the background. It will notify you once an appointment has been successfully booked.

---

## 🐛 Debugging

If you want to see what the bot is doing in real-time, you can run it in debug mode. This will disable headless mode (the browser window will be visible) and enable verbose logging.

To run in debug mode, set the `DEBUG` environment variable to `true`:

```bash
DEBUG=true go run ./cmd/zulassungsstellebot
```

---

## ⚙️ Advanced Configuration

For advanced users, the navigation flow of the bot can be customized by editing the `configs/menu.json` file. This file defines the menu structure and the corresponding selectors that the bot uses to navigate to the appointment calendar.
