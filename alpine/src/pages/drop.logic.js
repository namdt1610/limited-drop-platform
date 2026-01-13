import { getDropStatus, purchaseDrop } from "../lib/api.js";
import {
  getProvinces,
  getDistricts,
  getWards,
} from "../lib/vietnam-addresses.js";
import {
  initRevealObserver,
  formatTime,
  formatCountdown,
} from "../lib/dom-utils.js";

export function DropLogic() {
  return {
    dropId: null,
    dropData: null,
    isLoading: true,
    phase: "WAITING",
    countdown: "",
    remaining: 0,
    isSoldOut: false,
    winner: null,
    fomoStatuses: [],
    fomoTick: 0,

    // Checkout modal
    isModalOpen: false,
    isPurchasing: false,
    contact: {
      phone: "",
      email: "",
      name: "",
      address: "",
      province: "",
      district: "",
      ward: "",
    },
    magneticStyle: "",

    // Vietnam addresses
    provinces: [],
    districts: [],
    wards: [],
    provincesLoading: false,
    selectedProvinceCode: "",
    selectedDistrictCode: "",
    selectedWardCode: "",
    errors: {},

    serverOffset: 0,

    get isDisabled() {
      return !(this.phase === "LIVE" && !this.isSoldOut);
    },

    init() {
      window.addEventListener("route-changed", (event) => {
        if (event.detail.route === "drop") {
          this.dropId = event.detail.params.id || "1";
          this.loadDropData();
        }
      });

      // Initial load from hash
      const params = window.getHashParams();
      this.dropId = params.id || "1";
      this.loadDropData();
      this.loadProvinces();

      setInterval(() => {
        this.fomoTick++;
        this.updateFomoStatuses();
      }, 1500);

      this.initRevealObserver();
    },

    async loadProvinces() {
      this.provincesLoading = true;
      try {
        this.provinces = await getProvinces();
      } finally {
        this.provincesLoading = false;
      }
    },

    async onProvinceChange(provinceCode) {
      this.selectedProvinceCode = provinceCode;
      const province = this.provinces.find((p) => p.value === provinceCode);
      this.contact.province = province ? province.name : "";

      this.selectedDistrictCode = "";
      this.contact.district = "";
      this.selectedWardCode = "";
      this.contact.ward = "";

      this.districts = await getDistricts(provinceCode);
      this.wards = [];
      this.validateField("province");
    },

    async onDistrictChange(districtCode) {
      this.selectedDistrictCode = districtCode;
      const district = this.districts.find((d) => d.value === districtCode);
      this.contact.district = district ? district.name : "";

      this.selectedWardCode = "";
      this.contact.ward = "";

      this.wards = await getWards(districtCode);
      this.validateField("district");
    },

    onWardChange(wardCode) {
      this.selectedWardCode = wardCode;
      const ward = this.wards.find((w) => w.value === wardCode);
      this.contact.ward = ward ? ward.name : "";
      this.validateField("ward");
    },

    validateField(field) {
      const val = this.contact[field];
      if (!val) {
        this.errors[field] = "Trường này là bắt buộc";
        return;
      }

      if (field === "phone") {
        if (!/^(0|\+84)[3|5|7|8|9][0-9]{8}$/.test(val)) {
          this.errors[field] = "Số điện thoại không hợp lệ";
        } else {
          delete this.errors[field];
        }
      } else if (field === "email") {
        if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(val)) {
          this.errors[field] = "Email không hợp lệ";
        } else {
          delete this.errors[field];
        }
      } else {
        delete this.errors[field];
      }
    },

    validate() {
      const fields = [
        "phone",
        "email",
        "name",
        "address",
        "province",
        "district",
        "ward",
      ];
      fields.forEach((f) => this.validateField(f));
      return Object.keys(this.errors).length === 0;
    },

    initRevealObserver() {
      initRevealObserver();
    },

    handleMagnetic(e) {
      const btn = e.currentTarget;
      const rect = btn.getBoundingClientRect();
      const x = e.clientX - rect.left - rect.width / 2;
      const y = e.clientY - rect.top - rect.height / 2;
      this.magneticStyle = `transform: translate(${x * 0.3}px, ${y * 0.3}px)`;
    },

    resetMagnetic() {
      this.magneticStyle =
        "transform: translate(0, 0); transition: transform 0.5s cubic-bezier(0.23, 1, 0.32, 1)";
    },

    async loadDropData() {
      try {
        this.dropData = await getDropStatus(this.dropId);

        // Calculate server time offset
        const serverNow = new Date(this.dropData.now);
        const clientNow = new Date();
        this.serverOffset = serverNow - clientNow;

        this.updatePhase();
        this.isLoading = false;
        this.startCountdown();
      } catch (err) {
        this.isLoading = false;
      }
    },

    updatePhase() {
      if (!this.dropData) return;

      const now = new Date(Date.now() + this.serverOffset);
      const start = new Date(this.dropData.starts_at);
      const end = this.dropData.ends_at
        ? new Date(this.dropData.ends_at)
        : null;

      this.isSoldOut = Number(this.dropData.available) <= 0;

      if (now < start) {
        this.phase = "WAITING";
      } else if (this.isSoldOut) {
        this.phase = "SOLD_OUT";
      } else if (end && now >= end) {
        this.phase = "ENDED";
      } else {
        this.phase = "LIVE";
      }
    },

    startCountdown() {
      const update = () => {
        this.updatePhase();

        const now = new Date(Date.now() + this.serverOffset);
        const start = new Date(this.dropData.starts_at);
        const end = this.dropData.ends_at
          ? new Date(this.dropData.ends_at)
          : null;

        if (this.phase === "WAITING") {
          const diff = start - now;
          this.countdown = this.formatTime(diff);
          this.remaining = Math.floor(diff / 1000);
        } else if (this.phase === "LIVE") {
          const diff = end - now;
          this.countdown = this.formatTime(diff);
          this.remaining = Math.floor(diff / 1000);
        } else {
          this.countdown = "00:00:00";
          this.remaining = 0;
        }
      };
      update();
      setInterval(update, 1000);
    },

    formatTime(ms) {
      return formatTime(ms);
    },

    updateFomoStatuses() {
      const base = [
        { phone: "09x xxx 888", action: "đang thanh toán", time: "0.4s" },
        { phone: "08x xxx 555", action: "vừa giữ slot", time: "0.7s" },
        { phone: "03x xxx 123", action: "đang nhập địa chỉ", time: "0.9s" },
        {
          phone: "07x xxx 246",
          action: "đang kiểm tra thông tin",
          time: "1.1s",
        },
      ];
      const user =
        this.contact.phone && this.phase === "LIVE"
          ? [
              {
                phone: this.maskPhone(this.contact.phone),
                action: "đang tranh slot",
                time: "...",
              },
            ]
          : [];
      const winner =
        this.isSoldOut && this.winner
          ? [
              {
                phone: this.winner.maskedPhone,
                action: "đã chốt",
                time: this.winner.time || "—",
              },
            ]
          : [];
      const all = [...winner, ...base, ...user];
      if (!all.length) {
        this.fomoStatuses = [];
        return;
      }
      const start = this.fomoTick % all.length;
      this.fomoStatuses = [...all.slice(start), ...all.slice(0, start)].slice(
        0,
        Math.min(4, all.length)
      );
    },

    maskPhone(p) {
      if (!p) return "";
      return p.replace(/(\d{3})\d{3}(\d{3})/, "$1***$2");
    },
    openModal() {
      if (this.phase !== "LIVE" || this.isSoldOut) return;
      this.isModalOpen = true;
    },
    closeModal() {
      this.isModalOpen = false;
    },

    async handlePurchase() {
      if (!this.validate()) return;
      this.isPurchasing = true;
      try {
        const payload = {
          quantity: 1,
          name: this.contact.name || "Customer",
          phone: this.contact.phone,
          email: this.contact.email,
          address: this.contact.address || "",
          province: this.contact.province || "",
          district: this.contact.district || "",
          ward: this.contact.ward || "",
        };
        const res = await purchaseDrop(this.dropId, payload);
        if (res.payment_url) window.location.href = res.payment_url;
        else throw new Error("No payment URL received");
      } catch (err) {
        alert(err.message || "Có lỗi xảy ra. Vui lòng thử lại.");
      } finally {
        this.isPurchasing = false;
      }
    },
  };
}
