// Payment Cancel Page (Web Component + Alpine)
import { PaymentCancelLogic } from "./payment-cancel.logic.js";

// Định nghĩa custom element: <payment-cancel-page>
customElements.define(
  "payment-cancel-page",
  class extends HTMLElement {
    connectedCallback() {
      // Register Alpine data when element is connected
      if (!window.paymentCancelDataRegistered) {
        Alpine.data("paymentCancel", PaymentCancelLogic);
        window.paymentCancelDataRegistered = true;
      }

      this.innerHTML = /*html*/ `
<div x-data="paymentCancel" x-init="init()" class="relative min-h-screen bg-black text-white font-inter selection:bg-white selection:text-black">
  <!-- Background -->
  <mercury-background class="fixed inset-0 z-0 opacity-40"></mercury-background>

  <!-- Header -->
  <header class="reveal py-8 px-6 md:px-12 relative z-10">
    <div class="max-w-screen-2xl mx-auto flex items-center justify-between">
      <a href="#landing" @click.prevent="$dispatch('route', 'landing')" class="group">
        <img src="/imgs/logo-png.webp" alt="DONALD" class="w-20 md:w-24 h-auto object-contain transition-transform duration-500 group-hover:scale-110" />
      </a>
      <div class="flex flex-col items-end">
        <span class="text-[9px] font-mono text-white/30 uppercase tracking-[0.4em] mb-1">Giao thức</span>
        <span class="text-xs font-bold uppercase tracking-widest">Giao dịch bị Hủy</span>
      </div>
    </div>
  </header>

  <!-- Main Content -->
  <main class="relative z-10 px-6 md:px-12 pb-32">
    <div class="max-w-2xl mx-auto">
      <div class="reveal py-12 md:py-20 space-y-12 text-center">
        <div class="w-24 h-24 mx-auto bg-red-500/10 rounded-full flex items-center justify-center border border-red-500/20">
          <svg class="w-10 h-10 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </div>
        
        <div class="space-y-6">
          <h1 class="text-5xl md:text-7xl font-spiky tracking-tighter leading-none uppercase">
            GIAO DỊCH<br/><span class="text-red-500">BỊ HỦY</span>
          </h1>
          <p class="text-lg text-white/60 leading-relaxed max-w-md mx-auto">
            Giao thức đã bị ngắt quãng. Không có khoản phí nào được thực hiện.
          </p>
        </div>

        <div class="p-8 md:p-12 rounded-[2.5rem] bg-zinc-900/50 border border-white/10 backdrop-blur-3xl space-y-10 text-left">
          <div class="flex items-center gap-4">
            <span class="text-[10px] font-mono text-white/20 uppercase tracking-[0.4em]">01</span>
            <h2 class="text-xs font-bold uppercase tracking-widest">Chi tiết Giao dịch</h2>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div class="space-y-1">
              <span class="text-[9px] font-mono text-white/20 uppercase tracking-[0.4em]">Mã đơn hàng</span>
              <p class="text-sm font-mono text-white/80" x-text="orderCode || 'N/A'"></p>
            </div>
            <div class="space-y-1">
              <span class="text-[9px] font-mono text-white/20 uppercase tracking-[0.4em]">Trạng thái</span>
              <p class="text-sm font-mono text-red-500 uppercase tracking-widest" x-text="isCancelled ? 'Đã hủy' : 'Không xác định'"></p>
            </div>
            <div class="space-y-1" x-show="dropId">
              <span class="text-[9px] font-mono text-white/20 uppercase tracking-[0.4em]">ID Drop</span>
              <p class="text-sm font-mono text-white/80" x-text="'#' + dropId"></p>
            </div>
            <div class="space-y-1" x-show="code">
              <span class="text-[9px] font-mono text-white/20 uppercase tracking-[0.4em]">Mã phản hồi</span>
              <p class="text-sm font-mono text-white/80" x-text="code"></p>
            </div>
          </div>

          <div class="pt-8 border-t border-white/5 space-y-6">
            <div class="flex flex-col md:flex-row gap-4">
              <button @click="goHome" class="flex-1 py-6 rounded-full bg-white text-black text-[10px] font-bold tracking-[0.4em] uppercase transition-all duration-500 hover:scale-[1.02]">
                Quay lại Lõi
              </button>
              <button x-show="dropId" @click="retryPayment" class="flex-1 py-6 rounded-full bg-white/5 border border-white/10 text-white text-[10px] font-bold tracking-[0.4em] uppercase transition-all duration-500 hover:bg-white/10">
                Thử lại Giao dịch
              </button>
            </div>
          </div>
        </div>

        <div class="space-y-4 text-center">
          <p class="text-[10px] font-mono text-white/20 uppercase tracking-[0.4em] leading-relaxed">
            "Nếu đây là một lỗi hệ thống, vui lòng liên hệ với chúng tôi.
            <br/>Giao thức có thể được khởi tạo lại bất cứ lúc nào."
          </p>
        </div>
      </div>
    </div>
  </main>

  <!-- Footer -->
  <footer class="reveal py-20 px-6 md:px-12 border-t border-white/5">
    <div class="max-w-screen-2xl mx-auto text-center">
      <div class="text-[9px] font-mono text-white/10 uppercase tracking-[0.4em]">
        © 2026 DONALD CLUB
      </div>
    </div>
  </footer>
</div>
`;
    }
  }
);
