// Teaser Story Section Component for Alpine.js
document.addEventListener("alpine:init", () => {
  Alpine.data("teaserStorySection", () => ({
    tiltX: 0,
    tiltY: 0,

    handleTilt(e) {
      const card = e.currentTarget;
      const rect = card.getBoundingClientRect();
      const x = e.clientX - rect.left;
      const y = e.clientY - rect.top;
      const centerX = rect.width / 2;
      const centerY = rect.height / 2;
      this.tiltX = (y - centerY) / 10;
      this.tiltY = (centerX - x) / 10;
    },

    resetTilt() {
      this.tiltX = 0;
      this.tiltY = 0;
    },
  }));
});

// Define <teaser-story-section> directly as a custom element (inline template)
customElements.define(
  "teaser-story-section",
  class extends HTMLElement {
    connectedCallback() {
      this.innerHTML = /*html*/ `
    <section x-data="teaserStorySection" class="relative py-16 sm:py-20 lg:py-24 px-6 sm:px-8 md:px-12 overflow-hidden">
        <!-- Background gradient - subtle parallax -->
        <div class="absolute inset-0 bg-gradient-to-b from-transparent via-black/20 to-transparent pointer-events-none"></div>

        <div class="max-w-6xl mx-auto relative z-10">
            <div class="chrome-card p-8 lg:p-12 transition-all duration-500 hover:shadow-[0_0_50px_rgba(255,255,255,0.05)]">
            <!-- Main Story Narrative -->
            <div class="text-center mb-12 lg:mb-16 space-y-6 reveal">
                <div class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-card/50 border theme-border backdrop-blur transition-transform hover:scale-105">
                    <span class="text-xs font-medium tracking-[0.14em] text-muted-foreground uppercase">
                        Câu Chuyện Phía Sau
                    </span>
                </div>

                <div class="max-w-4xl mx-auto">
                    <h2 class="text-3xl sm:text-4xl lg:text-5xl font-bold text-foreground leading-tight transition-all duration-700 hover:tracking-tight">
                        Sự Cộng Sinh Giữa
                        <span class="block text-muted-foreground">Công Nghệ & Nghệ Thuật</span>
                    </h2>
                    <p class="text-lg sm:text-xl text-muted-foreground leading-relaxed max-w-3xl mx-auto mt-6">
                        Mỗi chiếc đồng hồ là một artifact sống - kết hợp giữa kỹ thuật chế tác truyền thống
                        và công nghệ tiên tiến. Sự hoàn hảo không phải là mục tiêu, mà là kết quả tất yếu.
                    </p>
                </div>
            </div>

            <!-- Story Points Grid -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-8 lg:gap-12">
                <!-- Point 1 -->
                <div class="group reveal" 
                     @mousemove="handleTilt($event)" 
                     @mouseleave="resetTilt"
                     :style="'transform: perspective(1000px) rotateX(' + tiltX + 'deg) rotateY(' + tiltY + 'deg)'">
                    <div class="feature-card transition-all duration-300 group-hover:border-primary/30 group-hover:bg-white/[0.02]">
                        <div class="flex items-start gap-4">
                            <div class="flex-shrink-0 w-12 h-12 bg-primary rounded-lg flex items-center justify-center group-hover:scale-110 group-hover:shadow-[0_0_20px_rgba(59,130,246,0.5)] transition-all">
                                <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z"></path>
                                </svg>
                            </div>
                            <div class="flex-1">
                                <h3 class="text-lg font-semibold text-foreground mb-2 group-hover:text-primary transition-colors">
                                    Hình Dáng Claw Symbiote
                                </h3>
                                <p class="text-muted-foreground leading-relaxed">
                                    Thiết kế lấy cảm hứng từ symbiote - những xúc tu organic vươn ra ôm sát cổ tay.
                                    Mỗi chi tiết được chế tác thủ công, tạo nên sự hoàn hảo về mặt thẩm mỹ.
                                </p>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Point 2 -->
                <div class="group reveal reveal-delay-1"
                     @mousemove="handleTilt($event)" 
                     @mouseleave="resetTilt"
                     :style="'transform: perspective(1000px) rotateX(' + tiltX + 'deg) rotateY(' + tiltY + 'deg)'">
                    <div class="feature-card transition-all duration-300 group-hover:border-primary/30 group-hover:bg-white/[0.02]">
                        <div class="flex items-start gap-4">
                            <div class="flex-shrink-0 w-12 h-12 bg-primary rounded-lg flex items-center justify-center group-hover:scale-110 group-hover:shadow-[0_0_20px_rgba(59,130,246,0.5)] transition-all">
                                <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
                                </svg>
                            </div>
                            <div class="flex-1">
                                <h3 class="text-lg font-semibold text-foreground mb-2 group-hover:text-primary transition-colors">
                                    Chất Liệu Chrome Living
                                </h3>
                                <p class="text-muted-foreground leading-relaxed">
                                    Bề mặt chrome sống động với hiệu ứng ánh sáng thay đổi theo góc nhìn.
                                    Sự phản chiếu hoàn hảo tạo nên vẻ ngoài alien và quyền lực.
                                </p>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Point 3 -->
                <div class="group reveal reveal-delay-2">
                    <div class="feature-card">
                        <div class="flex items-start gap-4">
                            <div class="flex-shrink-0 w-12 h-12 bg-primary rounded-lg flex items-center justify-center group-hover:scale-110 transition-transform">
                                <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"></path>
                                </svg>
                            </div>
                            <div class="flex-1">
                                <h3 class="text-lg font-semibold text-foreground mb-2 group-hover:text-primary transition-colors">
                                    Tính Cách Symbiote
                                </h3>
                                <p class="text-muted-foreground leading-relaxed">
                                    Mỗi artifact là duy nhất với pattern symbiote tự nhiên. Không có hai chiếc giống nhau - mỗi cái đều có DNA riêng của mình.
                                </p>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Point 4 -->
                <div class="group">
                    <div class="feature-card">
                        <div class="flex items-start gap-4">
                            <div class="flex-shrink-0 w-12 h-12 bg-primary rounded-lg flex items-center justify-center group-hover:scale-110 transition-transform">
                                <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 3l14 9-14 9V3z"></path>
                                </svg>
                            </div>
                            <div class="flex-1">
                                <h3 class="text-lg font-semibold text-foreground mb-2 group-hover:text-primary transition-colors">
                                    Sự Xứng Đáng
                                </h3>
                                <p class="text-muted-foreground leading-relaxed">
                                    Artifact chỉ chọn những vật chủ xứng đáng. Sự cộng sinh không phải may rủi - mà là sự nhận diện giữa tâm hồn tương đồng. Bạn có xứng đáng để trở thành một phần của câu chuyện này?
                                </p>
                            </div>
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
