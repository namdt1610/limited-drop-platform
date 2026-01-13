// Tech Specs Section Component for Alpine.js
document.addEventListener("alpine:init", () => {
  Alpine.data("techSpecsSection", () => ({
    activeSpec: 0,
    specs: [
      {
        title: "Vỏ Thép 904L",
        description:
          "Thép không gỉ cấp độ hàng không vũ trụ, được đánh bóng chrome gương thủ công.",
        icon: "M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z",
      },
      {
        title: "Kính Sapphire",
        description:
          "Chống trầy xước tuyệt đối với lớp phủ AR 7 lớp giảm phản xạ ánh sáng.",
        icon: "M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9",
      },
      {
        title: "Bộ Máy Automatic",
        description:
          "Tần số 28,800 vph, dự trữ năng lượng 72 giờ với rotor mạ rhodium.",
        icon: "M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z",
      },
      {
        title: "Kháng Nước 10ATM",
        description:
          "Hệ thống gioăng cao su kép đảm bảo an toàn ở độ sâu 100 mét.",
        icon: "M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.673.337a4 4 0 01-2.574.345l-2.387-.477a2 2 0 00-1.022.547l-1.162 1.162a2 2 0 01-2.828 0l-.141-.141a2 2 0 010-2.828l1.162-1.162a2 2 0 00.547-1.022l.477-2.387a6 6 0 00-.517-3.86l-.337-.673a4 4 0 01-.345-2.574l.477-2.387a2 2 0 00-.547-1.022l-1.162-1.162a2 2 0 010-2.828l.141-.141a2 2 0 012.828 0l1.162 1.162a2 2 0 001.022.547l2.387.477a6 6 0 003.86-.517l.673-.337a4 4 0 012.574-.345l2.387.477a2 2 0 001.022-.547l1.162-1.162a2 2 0 012.828 0l.141.141a2 2 0 010 2.828l-1.162 1.162a2 2 0 00-.547 1.022l-.477 2.387a6 6 0 00.517 3.86l.337.673a4 4 0 01.345 2.574l-.477 2.387a2 2 0 00.547 1.022l1.162 1.162a2 2 0 010 2.828l-.141.141a2 2 0 01-2.828 0l-1.162-1.162z",
      },
    ],
  }));
});

customElements.define(
  "tech-specs-section",
  class extends HTMLElement {
    connectedCallback() {
      this.innerHTML = /*html*/ `
        <section x-data="techSpecsSection" class="relative py-24 px-6 sm:px-8 md:px-12 overflow-hidden bg-black/40 backdrop-blur-sm border-y border-white/5">
            <div class="max-w-7xl mx-auto">
                <div class="grid grid-cols-1 lg:grid-cols-2 gap-16 items-center">
                    <!-- Left: Interactive List -->
                    <div class="space-y-8 reveal">
                        <div class="space-y-4">
                            <h2 class="text-3xl sm:text-4xl font-black tracking-tighter text-white uppercase">
                                THÔNG SỐ<br/><span class="text-white/40">KỸ THUẬT</span>
                            </h2>
                            <p class="text-zinc-500 font-mono text-sm uppercase tracking-widest">Protocol: Symbiosis-Specs-v1</p>
                        </div>

                        <div class="space-y-4">
                            <template x-for="(spec, index) in specs" :key="index">
                                <div 
                                    @mouseenter="activeSpec = index"
                                    class="group cursor-pointer p-6 rounded-2xl border transition-all duration-500"
                                    :class="activeSpec === index ? 'bg-white/5 border-white/20 shadow-2xl translate-x-4' : 'border-transparent hover:border-white/10'"
                                >
                                    <div class="flex items-center gap-6">
                                        <div 
                                            class="w-12 h-12 rounded-xl flex items-center justify-center transition-all duration-500"
                                            :class="activeSpec === index ? 'bg-primary text-white scale-110' : 'bg-white/5 text-zinc-600'"
                                        >
                                            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="spec.icon"></path>
                                            </svg>
                                        </div>
                                        <div class="flex-1">
                                            <h3 
                                                class="text-lg font-bold uppercase tracking-tight transition-colors duration-500"
                                                :class="activeSpec === index ? 'text-white' : 'text-zinc-500'"
                                                x-text="spec.title"
                                            ></h3>
                                            <p 
                                                class="text-sm leading-relaxed transition-all duration-500 overflow-hidden"
                                                :class="activeSpec === index ? 'text-zinc-400 max-h-20 mt-2 opacity-100' : 'text-zinc-600 max-h-0 opacity-0'"
                                                x-text="spec.description"
                                            ></p>
                                        </div>
                                        <div 
                                            class="transition-all duration-500"
                                            :class="activeSpec === index ? 'opacity-100 translate-x-0' : 'opacity-0 -translate-x-4'"
                                        >
                                            <svg class="w-5 h-5 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
                                            </svg>
                                        </div>
                                    </div>
                                </div>
                            </template>
                        </div>
                    </div>

                    <!-- Right: Visual Representation -->
                    <div class="relative aspect-square reveal reveal-delay-2">
                        <div class="absolute inset-0 bg-primary/10 rounded-full blur-[120px] animate-pulse"></div>
                        <div class="relative w-full h-full flex items-center justify-center">
                            <!-- Rotating Tech Ring -->
                            <div class="absolute inset-0 border-2 border-dashed border-white/5 rounded-full animate-[spin_20s_linear_infinite]"></div>
                            <div class="absolute inset-12 border border-primary/20 rounded-full animate-[spin_15s_linear_infinite_reverse]"></div>
                            
                            <!-- Active Spec Detail Card -->
                            <div class="relative z-10 card--glass p-8 rounded-3xl border border-white/10 shadow-2xl max-w-xs text-center space-y-6 transform transition-all duration-700"
                                 :class="activeSpec !== null ? 'scale-100 opacity-100' : 'scale-90 opacity-0'">
                                <div class="w-20 h-20 mx-auto bg-primary/20 rounded-2xl flex items-center justify-center text-primary">
                                    <svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="specs[activeSpec].icon"></path>
                                    </svg>
                                </div>
                                <div class="space-y-2">
                                    <h4 class="text-xl font-black text-white uppercase tracking-tighter" x-text="specs[activeSpec].title"></h4>
                                    <p class="text-zinc-400 text-sm leading-relaxed" x-text="specs[activeSpec].description"></p>
                                </div>
                                <div class="pt-4 border-t border-white/5">
                                    <span class="text-[10px] font-mono text-primary uppercase tracking-[0.3em]">Verified Component</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>
      `;
    }
  }
);
