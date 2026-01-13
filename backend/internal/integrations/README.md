# Integrations

This package contains code that interfaces with external systems and infra: Google APIs, Cloudinary, payment providers, email/resend services, and validation helpers.

Purpose:
- Keep third-party API wrappers and infra-related helpers isolated from core business logic.
- Provide simple, reusable functions for the rest of the application.

Files in this folder should expose minimal, well-documented functions and **not** contain domain business orchestration (those belong to `internal/service`).
