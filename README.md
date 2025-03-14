# go-url-shortener
Overview

This is a high-performance, secure URL shortener built with Golang, PostgreSQL, Redis, and JWT authentication. It allows users to shorten URLs, create custom aliases, track analytics, and ensure secure API access.

Features:

✅ Core Functionality

-> Shorten URLs – Convert long URLs into short links.
-> Custom Aliases – Create personalized short links.
-> Redirect Users – Automatically redirect visitors to the original URL.
-> Security
-> JWT Authentication – Secure access to shortening and analytics APIs.
-> Rate Limiting – Prevent API abuse and DDoS attacks.
-> Spam & Malicious URL Blocking – Protect against harmful sites.
-> Analytics
-> Track Clicks – Count the number of visits for each short URL.
-> View Last Accessed Time – Track when a short URL was last used.

-> Redis Caching – Improve performance with fast lookups.
-> URL Expiry & Cleanup – Automatically delete expired URLs.

🛠️ Tech Stack
Golang
Fiber
PostgreSQL
Redis
Caching & rate limiting
JWT Authentication

