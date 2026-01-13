// Google Form service - Headless Google Form submit helper
export async function submitToWaitlist({ email, phone } = {}) {
  const FORM_URL = import.meta.env.VITE_GOOGLE_FORM_URL;
  const EMAIL_ENTRY = import.meta.env.VITE_GOOGLE_FORM_EMAIL_ENTRY;
  const PHONE_ENTRY = import.meta.env.VITE_GOOGLE_FORM_PHONE_ENTRY;

  if (!FORM_URL || !EMAIL_ENTRY) {
    if (import.meta.env.DEV) {
      return;
    }
    throw new Error("Waitlist chưa được cấu hình");
  }

  const formData = { [EMAIL_ENTRY]: email };
  if (phone && PHONE_ENTRY) formData[PHONE_ENTRY] = phone;

  await fetch(FORM_URL, {
    method: "POST",
    mode: "no-cors",
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
    body: new URLSearchParams(formData),
  });
}
