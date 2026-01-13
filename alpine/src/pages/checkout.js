// Checkout Page (Web Component + Alpine)
import { CheckoutLogic } from "./checkout.logic.js";

// Định nghĩa custom element: <checkout-page>
customElements.define(
  "checkout-page",
  class extends HTMLElement {
    connectedCallback() {
      this.innerHTML = /*html*/ `
<div x-data="checkoutPage" class="relative min-h-screen bg-black text-white font-inter selection:bg-white selection:text-black">
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
        <span class="text-xs font-bold uppercase tracking-widest">Hoàn tất Giao dịch</span>
      </div>
    </div>
  </header>

  <!-- Main Content -->
  <main class="relative z-10 px-6 md:px-12 pb-32">
    <div class="max-w-screen-2xl mx-auto">
      <!-- Hero Section -->
      <div class="reveal py-12 md:py-20 space-y-6 text-center">
        <div class="inline-flex items-center gap-3 px-4 py-2 rounded-full bg-white/5 border border-white/10 backdrop-blur-xl">
          <span class="w-1.5 h-1.5 rounded-full bg-blue-500 animate-pulse"></span>
          <span class="text-[10px] font-mono text-white/60 uppercase tracking-[0.3em]">Kênh Bảo mật Đang hoạt động</span>
        </div>
        <h1 class="text-5xl md:text-7xl font-spiky tracking-tighter leading-none uppercase">
          HOÀN TẤT<br/><span class="text-white/20">GIAO DỊCH</span>
        </h1>
      </div>

      <div x-show="cart.length === 0" class="reveal py-32 text-center space-y-8">
        <p class="text-xs font-mono text-white/20 uppercase tracking-[0.4em]">Phát hiện Khoang trống</p>
        <a href="#products" @click.prevent="$dispatch('route', 'products')" class="inline-block px-12 py-5 rounded-full bg-white text-black text-[10px] font-bold uppercase tracking-[0.4em] hover:scale-105 transition-transform">
          Quay lại Lưu trữ
        </a>
      </div>

      <template x-if="cart.length > 0">
        <div class="grid grid-cols-1 lg:grid-cols-1 gap-12 lg:gap-20 items-start">
          <!-- Order Summary -->
          <div class="order-2 lg:order-1 space-y-6 sm:space-y-8 reveal px-0 sm:px-4" style="transition-delay: 200ms">
            <div class="p-6 sm:p-8 rounded-2xl sm:rounded-3xl bg-white/[0.02] border border-white/5 backdrop-blur-3xl space-y-6 sm:space-y-8">
              <div class="flex items-center gap-4">
                <span class="text-[9px] sm:text-[10px] font-mono text-white/20 uppercase tracking-[0.4em]">01</span>
                <h2 class="text-[11px] sm:text-xs font-bold uppercase tracking-widest">Tóm tắt Đơn hàng</h2>
              </div>
              
              <div class="space-y-4 sm:space-y-6">
                <template x-for="item in cart" :key="item.id">
                  <div class="flex justify-between items-center py-3 sm:py-4 border-b border-white/5">
                    <div class="space-y-1">
                      <p class="text-xs sm:text-sm font-bold uppercase tracking-tight" x-text="item.name"></p>
                      <p class="text-[8px] sm:text-[9px] font-mono text-white/20 uppercase">Số lượng: <span class="text-white/60" x-text="item.qty"></span></p>
                    </div>
                    <p class="text-xs sm:text-sm font-mono text-white/60" x-text="new Intl.NumberFormat('vi-VN').format(item.price * item.qty) + ' ₫'"></p>
                  </div>
                </template>
              </div>

              <div class="pt-3 sm:pt-4 flex justify-between items-end">
                <span class="text-[8px] sm:text-[9px] font-mono text-white/20 uppercase tracking-[0.4em]">Tổng giá trị Giao thức</span>
                <p class="text-2xl sm:text-3xl font-spiky tracking-tighter" x-text="new Intl.NumberFormat('vi-VN').format(getTotal()) + ' ₫'"></p>
              </div>
            </div>

            <div class="p-4 sm:p-6 rounded-xl sm:rounded-2xl bg-blue-500/[0.02] border border-blue-500/10 text-center">
              <p class="text-[8px] sm:text-[9px] font-mono text-blue-400/40 uppercase tracking-[0.2em] leading-relaxed">
                "Mọi giao dịch đều được mã hóa đầu cuối. Danh tính của bạn được bảo vệ bởi Giao thức Cộng sinh."
              </p>
            </div>
          </div>

          <!-- Shipping Form -->
          <div class="order-1 lg:order-2 reveal px-0 sm:px-4" style="transition-delay: 400ms">
            <div class="p-6 sm:p-8 md:p-12 rounded-[2rem] sm:rounded-[2.5rem] bg-zinc-900/50 border border-white/10 backdrop-blur-3xl space-y-6 sm:space-y-8 md:space-y-10">
              <div class="flex items-center gap-4">
                <span class="text-[9px] sm:text-[10px] font-mono text-white/20 uppercase tracking-[0.4em]">02</span>
                <h2 class="text-[11px] sm:text-xs font-bold uppercase tracking-widest">Danh tính & Tọa độ</h2>
              </div>

              <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6 md:gap-8">
                <div class="space-y-2">
                  <label class="text-[8px] sm:text-[9px] font-mono text-white/30 uppercase tracking-[0.3em] ml-1">Họ và Tên</label>
                  <input x-model="contact.name" class="w-full bg-white/[0.03] border border-white/10 rounded-2xl px-4 sm:px-6 py-3 sm:py-5 text-sm text-white placeholder:text-white/10 focus:outline-none focus:border-white/40 transition-all" placeholder="Tên Thực thể" />
                </div>
                <div class="space-y-2">
                  <label class="text-[8px] sm:text-[9px] font-mono text-white/30 uppercase tracking-[0.3em] ml-1">Số Điện thoại</label>
                  <input x-model="contact.phone" class="w-full bg-white/[0.03] border border-white/10 rounded-2xl px-4 sm:px-6 py-3 sm:py-5 text-sm font-mono text-white placeholder:text-white/10 focus:outline-none focus:border-white/40 transition-all" placeholder="+84 ..." />
                </div>
                <div class="space-y-2">
                  <label class="text-[8px] sm:text-[9px] font-mono text-white/30 uppercase tracking-[0.3em] ml-1">Email Định danh</label>
                  <input x-model="contact.email" class="w-full bg-white/[0.03] border border-white/10 rounded-2xl px-4 sm:px-6 py-3 sm:py-5 text-sm font-mono text-white placeholder:text-white/10 focus:outline-none focus:border-white/40 transition-all" placeholder="identity@protocol.xyz" />
                </div>
                <div class="space-y-2">
                  <label class="text-[8px] sm:text-[9px] font-mono text-white/30 uppercase tracking-[0.3em] ml-1">Tọa độ Giao hàng</label>
                  <textarea x-model="contact.address" class="w-full bg-white/[0.03] border border-white/10 rounded-2xl px-4 sm:px-6 py-3 sm:py-5 text-sm text-white placeholder:text-white/10 focus:outline-none focus:border-white/40 transition-all min-h-[100px] sm:min-h-[120px]" placeholder="Đường, Quận, Thành phố..."></textarea>
                </div>
              </div>

              <div class="pt-6 sm:pt-8 border-t border-white/5 space-y-4 sm:space-y-6">
                <button @click="submit()" :disabled="isSubmitting" 
                        class="w-full py-4 sm:py-6 rounded-full bg-white text-black text-[9px] sm:text-[10px] font-bold tracking-[0.4em] uppercase transition-all duration-500 hover:scale-[1.02] disabled:opacity-20">
                  <span x-show="!isSubmitting">Khởi tạo Giao dịch</span>
                  <span x-show="isSubmitting" class="animate-pulse">Đang xử lý Giao thức...</span>
                </button>
                <p x-show="error" class="text-red-500 text-[8px] sm:text-[10px] text-center font-mono uppercase tracking-widest" x-text="error"></p>
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </main>

  <!-- Footer -->
  <footer class="reveal py-20 px-6 md:px-12 border-t border-white/5">
    <div class="max-w-screen-2xl mx-auto flex flex-col md:flex-row items-center justify-between gap-8">
      <div class="flex items-center gap-4">
        <span class="text-[9px] font-mono text-white/20 uppercase tracking-[0.4em]">Phiên bản Giao thức</span>
        <span class="text-[9px] font-mono text-white/40 uppercase tracking-[0.4em]">v1.0.4-SECURE</span>
      </div>
      <div class="text-[9px] font-mono text-white/20 uppercase tracking-[0.4em]">
        © 2026 DONALD CLUB
      </div>
    </div>
  </footer>
</div>
    `;

      // Đăng ký logic tách riêng (guarded)
      if (!Alpine.store("checkoutPageInitialized")) {
        Alpine.data("checkoutPage", CheckoutLogic);
        Alpine.store("checkoutPageInitialized", true);
      }
    }
  }
);
