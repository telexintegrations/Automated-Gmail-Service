# 📧 Email Auto-Reply API

This API automatically sends an email reply when a message is received. It is designed to work with **Gmail** and **Outlook** users, making integration simple by handling SMTP configuration internally.

## 🚀 Features

- 📩 **Auto-reply emails** to incoming messages
- 🔄 **Preconfigured SMTP settings** for Gmail and Outlook
- 🔐 **Supports Two-Factor Authentication (2FA) users** via App Passwords
- 🌍 **Cross-Origin Resource Sharing (CORS) support**

---

## ⚙️ API Setup

### **1️⃣ Install Dependencies**

Ensure you have Go installed, then run:

```
go mod tidy
```

### **2️⃣ Start the API**

Run the server with:

```
go run main.go
```

## **🔗 API Endpoints**

###### 📨Send Auto-Reply

Endpoint:

```
POST /start-email-listener
```

##### Request Body Example:

```
{
  "username": "Your Name",
  "email": "your-email@gmail.com",
  "password": "your-app-password"
}
```

#### 🔐 Using the API with 2FA-Enabled Accounts

We advise that don't use your regular email password. Instead, on your device navigate into Gmail and click on 'Manage your Google Account'. Locate the Security tab and click on it. Ensure you have enabled 2-step verification and generate an App Password from your email provider. This is what you'd use as your-app-password.

#### 🛠️ What is an App Password?

An App Password is a special password generated from your email provider’s security settings. It allows third-party applications (like this API) to send emails without disabling 2FA.

#### 🔑 How to Generate an App Password

Follow the steps below to create an App Password for Gmail and Outlook.

#### For Gmail Users

```
Go to Google App Passwords.
Select the app ("Mail") and device ("Other (Custom name)").
Click Generate – Google will provide you with a 16-character password.
Use this password instead of your real email password when making API requests.
```

#### For Outlook Users

```
Go to Microsoft Security Settings.
Click Create a new app password.
Copy the password and use it for your API requests.
```

#### A successful authentication response would look like:

```
{
  "details":"Email monitoring started. New inbox mails would receive automated responses.",
  "message":"Login successful."
}
```
