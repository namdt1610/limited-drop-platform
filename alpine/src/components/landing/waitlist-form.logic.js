// Logic for Waitlist Form component
// import { submitToWaitlist } from "../lib/google-form.service.js";

export const MESSAGES = {
  REQUIRED_EMAIL: "VẬT CHỦ ơi, VUI LÒNG CẬP NHẬT TẦN SỐ QUÉT CỦA BẠN",
  INVALID_EMAIL: "TẦN SỐ QUÉT KHÔNG HỢP LỆ (VẬT CHỦ CẦU THẬN)",
  REQUIRED_PHONE: "VẬT CHỦ, VUI LÒNG NHẬP KÊNH LIÊN LẠC CỦA BẠN",
  INVALID_PHONE:
    "KÊNH LIÊN LẠC KHÔNG HỢP LỆ (VD: 0901234567, +84901234567, 090-123-4567)",
  SUCCESS: "CHẤP NHẬN CÔNG SINH HOÀN THÀNH. CHUẨN BỊ ĐÓN VẬT CHỦ MỚI.",
  ERROR: "Có lỗi xảy ra. Vui lòng thử lại sau.",
};

export function validateEmail(email) {
  // More strict email validation
  const re =
    /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/;
  return re.test(
    String(email || "")
      .trim()
      .toLowerCase()
  );
}

export function validatePhone(phone) {
  if (!phone || !phone.trim()) return false;

  // Remove all non-digit characters except + at the beginning
  let cleaned = phone.trim().replace(/[^\d+]/g, "");

  // Handle international format +84
  if (cleaned.startsWith("+84")) {
    cleaned = "0" + cleaned.substring(3);
  }

  // Remove + if it's not at the beginning
  cleaned = cleaned.replace(/^\+/, "");

  // Vietnam mobile phone patterns (10-11 digits)
  const patterns = [
    /^0[3|5|7|8|9]\d{8}$/, // Mobile: 03x, 05x, 07x, 08x, 09x (10 digits)
    /^0[2]\d{9}$/, // Hanoi: 02x (10 digits, but now 11 digits)
    /^0[2]\d{10}$/, // Hanoi: 02x (11 digits)
  ];

  return patterns.some((pattern) => pattern.test(cleaned));
}

export async function submitWaitlist(payload) {
  // For Alpine-only development, simulate successful submission
  // TODO: Replace with actual API call when backend is ready
  console.log("[DEV] Waitlist submission:", payload);

  // Simulate network delay
  await new Promise((resolve) => setTimeout(resolve, 1000));

  return { success: true };
}
