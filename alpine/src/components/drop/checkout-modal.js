export function CheckoutModal() {
  return /*html*/ `
    <div x-show="isModalOpen" @keydown.escape.window="closeModal" class="fixed inset-0 z-[100] flex items-center justify-center p-4 sm:p-6 overflow-y-auto"
         x-transition:enter="transition ease-out duration-500" x-transition:enter-start="opacity-0" x-transition:enter-end="opacity-100"
         x-transition:leave="transition ease-in duration-300" x-transition:leave-start="opacity-100" x-transition:leave-end="opacity-0">
      
      <!-- Backdrop -->
      <div class="absolute inset-0 bg-black/90 backdrop-blur-md" @click="closeModal"></div>

      <!-- Modal Content -->
      <div class="relative w-full max-w-sm sm:max-w-2xl md:max-w-3xl lg:max-w-5xl bg-zinc-900/80 border border-white/10 p-6 sm:p-8 md:p-12 rounded-[2rem] backdrop-blur-3xl shadow-2xl space-y-6 sm:space-y-8 my-auto max-h-[90vh] overflow-y-auto"
           x-transition:enter="transition ease-out duration-500" x-transition:enter-start="opacity-0 scale-95 translate-y-8" x-transition:enter-end="opacity-100 scale-100 translate-y-0"
           x-transition:leave="transition ease-in duration-300" x-transition:leave-start="opacity-100 scale-100 translate-y-0" x-transition:leave-end="opacity-0 scale-95 translate-y-8">
        
        <!-- Header -->
        <div class="space-y-3 sm:space-y-4 text-center">
          <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-white/5 border border-white/10 text-[8px] sm:text-[9px] uppercase tracking-[0.3em] text-white/40 font-mono">
              Yêu cầu Xác minh Định danh
          </div>
          <h3 class="text-2xl sm:text-3xl md:text-4xl font-spiky tracking-tighter text-white uppercase leading-none">
              GIỮ<br/><span class="text-white/20">SLOT CỦA BẠN</span>
          </h3>
        </div>

        <form @submit.prevent="handlePurchase" class="space-y-5 sm:space-y-6 md:space-y-8">
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6 md:gap-8">
              <!-- Phone Input -->
              <div class="form-wrapper">
                  <label class="form-label">Số điện thoại</label>
                  <input 
                      x-model="contact.phone" 
                      @input="validateField('phone')"
                      @blur="validateField('phone')"
                      type="tel" 
                      placeholder="+84 ..." 
                      class="form-input"
                      :class="errors.phone ? 'error' : ''"
                      required 
                  />
                  <p x-show="errors.phone" x-text="errors.phone" class="form-error"></p>
              </div>

              <!-- Email Input -->
              <div class="form-wrapper">
                  <label class="form-label">Email Định danh</label>
                  <input 
                      x-model="contact.email" 
                      @input="validateField('email')"
                      @blur="validateField('email')"
                      type="email" 
                      placeholder="identity@protocol.xyz" 
                      class="form-input"
                      :class="errors.email ? 'error' : ''"
                      required 
                  />
                  <p x-show="errors.email" x-text="errors.email" class="form-error"></p>
              </div>

              <!-- Name Input -->
              <div class="form-wrapper">
                  <label class="form-label">Họ và Tên</label>
                  <input 
                      x-model="contact.name" 
                      @input="validateField('name')"
                      @blur="validateField('name')"
                      type="text" 
                      placeholder="Tên của bạn" 
                      class="form-input"
                      :class="errors.name ? 'error' : ''"
                  />
                  <p x-show="errors.name" x-text="errors.name" class="form-error"></p>
              </div>

              <!-- Address Input -->
              <div class="form-wrapper lg:col-span-3">
                  <label class="form-label">Địa chỉ</label>
                  <input 
                      x-model="contact.address" 
                      @input="validateField('address')"
                      @blur="validateField('address')"
                      type="text" 
                      placeholder="Địa chỉ giao hàng" 
                      class="form-input"
                      :class="errors.address ? 'error' : ''"
                  />
                  <p x-show="errors.address" x-text="errors.address" class="form-error"></p>
              </div>

              <!-- Province Select -->
              <div class="form-wrapper">
                  <label class="form-label">Tỉnh / Thành phố</label>
                  <div class="relative">
                      <select 
                          @change="onProvinceChange($event.target.value)"
                          :value="selectedProvinceCode"
                          :disabled="provincesLoading || provinces.length === 0"
                          class="form-select"
                          :class="errors.province ? 'error' : ''"
                      >
                          <option value="" disabled selected class="bg-zinc-900 text-white">Chọn tỉnh / thành phố</option>
                          <template x-for="province in provinces" :key="province.value">
                              <option :value="province.value" x-text="province.label" class="bg-zinc-900 text-white"></option>
                          </template>
                      </select>
                      <div class="absolute right-5 top-1/2 -translate-y-1/2 pointer-events-none text-white/20">
                          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
                      </div>
                  </div>
                  <p x-show="errors.province" x-text="errors.province" class="form-error"></p>
              </div>

              <!-- District Select -->
              <div class="form-wrapper">
                  <label class="form-label">Quận / Huyện</label>
                  <div class="relative">
                      <select 
                          @change="onDistrictChange($event.target.value)"
                          :value="selectedDistrictCode"
                          :disabled="!selectedProvinceCode || districts.length === 0"
                          class="form-select"
                          :class="errors.district ? 'error' : ''"
                      >
                          <option value="" disabled selected class="bg-zinc-900 text-white">Chọn quận / huyện</option>
                          <template x-for="district in districts" :key="district.value">
                              <option :value="district.value" x-text="district.label" class="bg-zinc-900 text-white"></option>
                          </template>
                      </select>
                      <div class="absolute right-5 top-1/2 -translate-y-1/2 pointer-events-none text-white/20">
                          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
                      </div>
                  </div>
                  <p x-show="errors.district" x-text="errors.district" class="form-error"></p>
              </div>

              <!-- Ward Select -->
              <div class="form-wrapper">
                  <label class="form-label">Phường / Xã</label>
                  <div class="relative">
                      <select 
                          @change="onWardChange($event.target.value)"
                          :value="selectedWardCode"
                          :disabled="!selectedDistrictCode || wards.length === 0"
                          class="form-select"
                          :class="errors.ward ? 'error' : ''"
                      >
                          <option value="" disabled selected class="bg-zinc-900 text-white">Chọn phường / xã</option>
                          <template x-for="ward in wards" :key="ward.value">
                              <option :value="ward.value" x-text="ward.label" class="bg-zinc-900 text-white"></option>
                          </template>
                      </select>
                      <div class="absolute right-5 top-1/2 -translate-y-1/2 pointer-events-none text-white/20">
                          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
                      </div>
                  </div>
                  <p x-show="errors.ward" x-text="errors.ward" class="form-error"></p>
              </div>
          </div>

          <div class="p-4 sm:p-5 rounded-2xl bg-white/[0.02] border border-white/5 text-center">
              <p class="text-[8px] sm:text-[9px] font-mono text-white/40 uppercase tracking-[0.2em] leading-relaxed">
                  "Bằng cách tiếp tục, bạn đồng ý với giao thức thanh toán 60 giây. Thất bại sẽ dẫn đến việc hủy bỏ slot."
              </p>
          </div>

          <div class="flex flex-col gap-3 sm:gap-4">
            <button type="submit" :disabled="isPurchasing || !contact.phone || !contact.email" 
                    class="w-full py-4 sm:py-6 rounded-full bg-white text-black text-[9px] sm:text-[10px] font-bold tracking-[0.4em] uppercase transition-all duration-500 hover:scale-[1.02] disabled:opacity-20 disabled:cursor-not-allowed">
              <span x-show="!isPurchasing">Khởi tạo Giao dịch</span>
              <span x-show="isPurchasing" class="animate-pulse">Đang xử lý Định danh...</span>
            </button>
            <button type="button" @click="closeModal" class="w-full py-3 sm:py-4 text-[8px] sm:text-[9px] font-mono text-white/20 uppercase tracking-[0.3em] hover:text-white/60 transition-colors">
              Hủy bỏ Giao thức
            </button>
          </div>
        </form>

        <div class="text-center">
          <p class="text-[7px] sm:text-[8px] font-mono text-white/10 uppercase tracking-[0.4em]">
              Giao thức Hệ thống: Symbiosis_Auth_v1.0.4
          </p>
        </div>
      </div>
    </div>
  `;
}
