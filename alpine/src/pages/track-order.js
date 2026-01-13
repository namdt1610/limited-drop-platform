// Track Order Page (Web Component + Alpine)
import { TrackOrderLogic } from "./track-order.logic.js";

// Định nghĩa custom element: <track-order-page>
customElements.define(
  "track-order-page",
  class extends HTMLElement {
    connectedCallback() {
      this.innerHTML = /*html*/ `
<div x-data="trackOrderPage" class="relative min-h-screen bg-black text-white font-inter selection:bg-white selection:text-black">
  <!-- Background -->
  <mercury-background class="fixed inset-0 z-0 opacity-40"></mercury-background>

  <!-- Header -->
  <header class="reveal py-8 px-6 md:px-12 relative z-10">
    <div class="max-w-screen-2xl mx-auto flex items-center justify-between">
      <a href="#landing" @click.prevent="$dispatch('route', 'landing')" class="group">
        <img src="/imgs/logo-png.webp" alt="DONALD" class="w-20 md:w-24 h-auto object-contain transition-transform duration-500 group-hover:scale-110" />
      </a>
      <div class="flex flex-col items-end">
        <span class="text-[9px] font-mono text-white/30 uppercase tracking-[0.4em] mb-1">Truy vấn</span>
        <span class="text-xs font-bold uppercase tracking-widest">Trạng thái Giao thức</span>
      </div>
    </div>
  </header>

  <!-- Main Content -->
  <main class="relative z-10 px-6 md:px-12 pb-32">
    <div class="max-w-4xl mx-auto">
      <!-- Hero Section -->
      <div class="reveal py-12 md:py-20 space-y-6 text-center">
        <div class="inline-flex items-center gap-3 px-4 py-2 rounded-full bg-white/5 border border-white/10 backdrop-blur-xl">
          <span class="w-1.5 h-1.5 rounded-full bg-blue-500 animate-pulse"></span>
          <span class="text-[10px] font-mono text-white/60 uppercase tracking-[0.3em]">Đang kết nối Cơ sở dữ liệu...</span>
        </div>
        <h1 class="text-5xl md:text-7xl font-spiky tracking-tighter leading-none uppercase">
          TRA CỨU<br/><span class="text-white/20">GIAO DỊCH</span>
        </h1>
      </div>

      <div class="reveal p-8 md:p-16 rounded-[2.5rem] bg-zinc-900/50 border border-white/10 backdrop-blur-3xl space-y-12">
        <template x-if="!order">
          <div class="space-y-10">
            <div class="space-y-6">
              <div class="flex items-center gap-4">
                <span class="text-[10px] font-mono text-white/20 uppercase tracking-[0.4em]">01</span>
                <h2 class="text-xs font-bold uppercase tracking-widest">Thông tin Truy vấn</h2>
              </div>
              
              <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
                <div class="space-y-2">
                  <label class="text-[9px] font-mono text-white/30 uppercase tracking-[0.3em] ml-1">Mã đơn hàng</label>
                  <input x-model="orderNumber" class="w-full bg-white/[0.03] border border-white/10 rounded-2xl px-6 py-5 text-sm font-mono text-white placeholder:text-white/10 focus:outline-none focus:border-white/40 transition-all" placeholder="DV-XXXXXX" />
                </div>
                <div class="space-y-2">
                  <label class="text-[9px] font-mono text-white/30 uppercase tracking-[0.3em] ml-1">Số điện thoại</label>
                  <input x-model="phone" class="w-full bg-white/[0.03] border border-white/10 rounded-2xl px-6 py-5 text-sm font-mono text-white placeholder:text-white/10 focus:outline-none focus:border-white/40 transition-all" placeholder="09xxxxxxxx" />
                </div>
              </div>
            </div>

            <button @click="submitForm" :disabled="loading" 
                    class="w-full py-6 rounded-full bg-white text-black text-[10px] font-bold tracking-[0.4em] uppercase transition-all duration-500 hover:scale-[1.02] disabled:opacity-20">
              <span x-show="!loading">Bắt đầu Truy vấn</span>
              <span x-show="loading" class="animate-pulse">Đang giải mã...</span>
            </button>
            
            <p x-show="error" class="text-red-500 text-[10px] text-center font-mono uppercase tracking-widest" x-text="error"></p>
          </div>
        </template>

        <template x-if="order">
          <div class="space-y-12">
            <!-- Order Status -->
            <div class="space-y-8">
              <div class="flex items-center justify-between">
                <div class="space-y-1">
                  <span class="text-[9px] font-mono text-white/20 uppercase tracking-[0.4em]">Trạng thái hiện tại</span>
                  <h3 class="text-2xl font-spiky tracking-tighter uppercase" :class="getStatusColor(order.status)" x-text="getStatusLabel(order.status)"></h3>
                </div>
                <div class="text-right space-y-1">
                  <span class="text-[9px] font-mono text-white/20 uppercase tracking-[0.4em]">Mã Giao dịch</span>
                  <p class="text-sm font-mono text-white/60" x-text="order.orderNumber"></p>
                </div>
              </div>

              <!-- Progress Steps -->
              <div class="grid grid-cols-4 gap-2">
                <template x-for="(step, index) in getStatusSteps(order.status)" :key="index">
                  <div class="space-y-3">
                    <div class="h-1 rounded-full transition-colors duration-1000" 
                         :class="step.completed ? 'bg-white' : (step.active ? 'bg-blue-500 animate-pulse' : 'bg-white/10')"></div>
                    <span class="block text-[8px] font-mono uppercase tracking-widest text-center" 
                          :class="step.active || step.completed ? 'text-white' : 'text-white/20'" x-text="step.label"></span>
                  </div>
                </template>
              </div>
            </div>

            <!-- Order Details -->
            <div class="pt-12 border-t border-white/5 grid grid-cols-1 md:grid-cols-2 gap-12">
              <div class="space-y-6">
                <h4 class="text-[10px] font-bold uppercase tracking-widest">Thông tin Thực thể</h4>
                <div class="space-y-4 text-sm font-light text-white/60">
                  <p x-text="order.shippingAddress.name"></p>
                  <p x-text="order.shippingAddress.phone"></p>
                  <p class="leading-relaxed" x-text="order.shippingAddress.detail + ', ' + order.shippingAddress.city"></p>
                </div>
              </div>
              <div class="space-y-6">
                <h4 class="text-[10px] font-bold uppercase tracking-widest">Tóm tắt Giao thức</h4>
                <div class="space-y-4">
                  <template x-for="item in order.items" :key="item.product_name">
                    <div class="flex justify-between text-sm">
                      <span class="text-white/40" x-text="item.quantity + 'x ' + item.product_name"></span>
                      <span class="font-mono" x-text="new Intl.NumberFormat('vi-VN').format(item.subtotal) + ' ₫'"></span>
                    </div>
                  </template>
                  <div class="pt-4 border-t border-white/5 flex justify-between items-end">
                    <span class="text-[9px] font-mono text-white/20 uppercase tracking-[0.4em]">Tổng giá trị</span>
                    <span class="text-xl font-spiky tracking-tighter" x-text="new Intl.NumberFormat('vi-VN').format(order.totalAmount) + ' ₫'"></span>
                  </div>
                </div>
              </div>
            </div>

            <div class="pt-8 text-center">
              <button @click="resetForm" class="text-[10px] font-bold uppercase tracking-[0.4em] text-white/40 hover:text-white transition-colors">
                Truy vấn mã khác
              </button>
            </div>
          </div>
        </template>
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

      // Đăng ký logic tách riêng (guarded)
      if (!Alpine.store("trackOrderPageInitialized")) {
        Alpine.data("trackOrderPage", TrackOrderLogic);
        Alpine.store("trackOrderPageInitialized", true);
      }
    }
  }
);
