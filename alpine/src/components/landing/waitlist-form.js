import {
  validateEmail,
  validatePhone,
  submitWaitlist,
  MESSAGES,
} from "./waitlist-form.logic.js";

class WaitlistForm extends HTMLElement {
  constructor() {
    super();
    // 1. Khởi tạo State với cơ chế "Phản ứng" (Reactive)
    this._state = {
      email: "",
      phone: "",
      message: "",
      isSubmitting: false,
      isSubscribed: false,
    };
  }

  // Cơ chế "Camera" theo dõi state
  set state(newState) {
    this._state = { ...this._state, ...newState };
    this.updateDOM(); // Tự động cập nhật chỉ những chỗ cần thiết
  }

  get state() {
    return this._state;
  }

  connectedCallback() {
    // 2. Chỉ vẽ khung (Shell) MỘT LẦN DUY NHẤT
    this.renderInitialHTML();
    // 3. Cache các phần tử DOM (để không phải query liên tục)
    this.cacheDOM();
    this.attachEventListeners();
  }

  cacheDOM() {
    this.form = this.querySelector("form");
    this.emailInput = this.querySelector('input[type="email"]');
    this.phoneInput = this.querySelector('input[type="tel"]');
    this.submitBtn = this.querySelector('button[type="submit"]');
    this.messageEl = this.querySelector(".status-message");
    this.btnText = this.querySelector(".btn-text");
    this.btnSpinner = this.querySelector(".btn-spinner");
  }

  // 4. KIẾN TRÚC REACTIVE: Chỉ thay đổi đúng cái Node cần thiết
  updateDOM() {
    const { message, isSubmitting, isSubscribed } = this.state;

    // Cập nhật text báo lỗi/thành công
    if (this.messageEl) {
      this.messageEl.textContent = message;
      this.messageEl.className = `status-message text-sm mb-2 transition-all ${
        message === MESSAGES.SUCCESS ? "text-green-400" : "text-red-400"
      }`;
    }

    // Cập nhật trạng thái nút bấm
    if (this.submitBtn) {
      const hasError = message && message !== MESSAGES.SUCCESS;
      this.submitBtn.disabled = isSubmitting || hasError;
      this.submitBtn.style.opacity = isSubmitting || hasError ? "0.5" : "1";
    }

    // Điều khiển Spinner
    this.btnSpinner.classList.toggle("hidden", !isSubmitting);
    this.btnText.textContent = isSubmitting ? "Đang gửi..." : "Nhận Thông Báo";

    // Nếu đã đăng ký thành công, mày có thể ẩn form/hiện success ở đây
    if (isSubscribed) {
      this.innerHTML = this.renderSuccess(); // Đoạn này mới cần thay cả cụm
    }
  }

  async handleSubmit(e) {
    e.preventDefault();
    const email = this.emailInput.value.trim();
    const phone = this.phoneInput.value.trim();

    // Validation logic (TDD ready)
    if (!validateEmail(email))
      return (this.state = { message: MESSAGES.INVALID_EMAIL });
    if (!validatePhone(phone))
      return (this.state = { message: MESSAGES.INVALID_PHONE });

    this.state = { isSubmitting: true, message: "" };

    try {
      await submitWaitlist({ email, phone });
      this.state = { isSubscribed: true, message: MESSAGES.SUCCESS };
    } catch (err) {
      this.state = { isSubmitting: false, message: MESSAGES.ERROR };
    }
  }

  attachEventListeners() {
    this.form.addEventListener("submit", (e) => this.handleSubmit(e));

    // Xóa lỗi khi người dùng bắt đầu sửa lại - Cực nhạy!
    this.emailInput.addEventListener("input", () => {
      if (this.state.message) this.state = { message: "" };
    });
    this.phoneInput.addEventListener("input", () => {
      if (this.state.message) this.state = { message: "" };
    });
  }

  renderInitialHTML() {
    this.innerHTML = /*html*/ `
      <form class="flex flex-col gap-3 max-w-md mx-auto lg:mx-0">
          <div class="status-message min-h-[20px]"></div>
          
          <div class="relative group">
              <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 12a4 4 0 10-8 0 4 4 0 008 0zm0 0v1.5a2.5 2.5 0 005 0V12a9 9 0 10-9 9m4.5-1.206a8.959 8.959 0 01-4.5 1.207"></path>
              </svg>
              <input class="pl-12 h-12 w-full bg-card border border-border rounded-md" type="email" placeholder="Email của bạn">
          </div>

          <div class="relative group">
              <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path>
              </svg>
              <input class="pl-12 h-12 w-full bg-card border border-border rounded-md" type="tel" placeholder="Số điện thoại">
          </div>

          <button class="chrome-btn w-full h-12 flex items-center justify-center rounded-md" type="submit">
              <span class="btn-spinner hidden w-4 h-4 border-2 border-t-transparent border-white rounded-full animate-spin mr-2"></span>
              <span class="btn-text">Nhận Thông Báo</span>
          </button>
      </form>
    `;
  }

  renderSuccess() {
    return `<div class="p-6 bg-gray-100/10 rounded-xl text-center">VẬT CHỦ ĐÃ SẴN SÀNG!</div>`;
  }
}

customElements.define("waitlist-form", WaitlistForm);
