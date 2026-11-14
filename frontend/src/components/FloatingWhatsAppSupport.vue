<script setup>
import { ref, onMounted, onBeforeUnmount } from "vue";
import {
  ExternalLink,
  Headphones,
  ShoppingCart,
  GraduationCap,
  X,
} from "lucide-vue-next";

const isOpen = ref(false);
const widgetRef = ref(null);

const contacts = [
  {
    id: 1,
    name: "CS Probussystem",
    role: "Customer Service & Bantuan Teknis",
    numberDisplay: "+62 819-1087-5514",
    waNumber: "6281910875514",
    badge: "CS Support",
    badgeClass: "bg-emerald-50 text-emerald-700 border-emerald-200",
    icon: Headphones,
    defaultMessage:
      "Halo CS Probussystem, saya butuh bantuan terkait kendala sistem Probus.",
  },
  {
    id: 2,
    name: "E Commerce Probussystem",
    role: "Online Booking, Channel & Integrasi",
    numberDisplay: "+62 878-6144-6124",
    waNumber: "6287861446124",
    badge: "E-Commerce",
    badgeClass: "bg-sky-50 text-sky-700 border-sky-200",
    icon: ShoppingCart,
    defaultMessage:
      "Halo E-Commerce Probussystem, saya ingin konsultasi terkait layanan integrasi & E-Commerce Probus.",
  },
  {
    id: 3,
    name: "Trainer Probussystem",
    role: "Panduan Operasional & Training Modul",
    numberDisplay: "+62 819-1087-5514",
    waNumber: "6281910875514",
    badge: "Trainer",
    badgeClass: "bg-indigo-50 text-indigo-700 border-indigo-200",
    icon: GraduationCap,
    defaultMessage:
      "Halo Trainer Probussystem, saya ingin panduan dan jadwal training modul sistem Probus.",
  },
];

function toggleOpen() {
  isOpen.value = !isOpen.value;
}

function openWhatsApp(contact) {
  const url = `https://wa.me/${contact.waNumber}?text=${encodeURIComponent(contact.defaultMessage)}`;
  window.open(url, "_blank");
  isOpen.value = false;
}

function handleClickOutside(event) {
  if (widgetRef.value && !widgetRef.value.contains(event.target)) {
    isOpen.value = false;
  }
}

function handleKeydown(event) {
  if (event.key === "Escape" && isOpen.value) {
    isOpen.value = false;
  }
}

onMounted(() => {
  document.addEventListener("click", handleClickOutside);
  document.addEventListener("keydown", handleKeydown);
});

onBeforeUnmount(() => {
  document.removeEventListener("click", handleClickOutside);
  document.removeEventListener("keydown", handleKeydown);
});
</script>

<template>
  <div
    class="fixed bottom-6 right-6 z-50 select-none font-sans"
    ref="widgetRef"
  >
    <!-- Popover Card Contact List (Opens Above the Button) -->
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="transform opacity-0 scale-95 translate-y-3"
      enter-to-class="transform opacity-100 scale-100 translate-y-0"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="transform opacity-100 scale-100 translate-y-0"
      leave-to-class="transform opacity-0 scale-95 translate-y-3"
    >
      <div
        v-if="isOpen"
        class="fixed bottom-24 right-4 sm:right-6 w-[calc(100vw-2rem)] sm:w-96 max-w-[390px] bg-white rounded-2xl border border-slate-200 shadow-2xl overflow-hidden z-50 text-left"
      >
        <!-- Header of Popup -->
        <div
          class="p-4 sm:p-4.5 bg-gradient-to-r from-emerald-600 via-emerald-600 to-teal-700 text-white"
        >
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2.5">
              <div
                class="w-8 h-8 rounded-xl bg-white/20 backdrop-blur-md flex items-center justify-center shrink-0 shadow-inner"
              >
                <!-- WhatsApp SVG Icon -->
                <svg viewBox="0 0 24 24" class="w-4.5 h-4.5 fill-white">
                  <path
                    d="M12.031 2C6.511 2 2.016 6.496 2.016 12.016c0 1.988.583 3.918 1.688 5.568L2.016 22l4.57-1.636c1.604.996 3.473 1.652 5.445 1.652 5.52 0 10.016-4.496 10.016-10.016 0-5.52-4.496-10-10.016-10zm.016 18.234c-1.745 0-3.414-.492-4.852-1.422l-.348-.225-2.871 1.028.98-2.784-.243-.377c-1.043-1.522-1.6-3.328-1.6-5.454 0-4.524 3.68-8.204 8.203-8.204 4.524 0 8.204 3.68 8.204 8.204 0 4.523-3.68 8.234-8.204 8.234zm4.496-6.148c-.246-.123-1.455-.718-1.68-.8-.225-.082-.389-.123-.553.123-.164.246-.635.8-.779.963-.143.164-.287.184-.533.061-.246-.123-1.037-.383-1.977-1.221-.73-.652-1.223-1.457-1.367-1.703-.143-.246-.015-.379.108-.502.111-.111.246-.287.369-.43.123-.144.164-.246.246-.41.082-.164.041-.308-.02-.43-.062-.124-.554-1.334-.759-1.826-.2-.48-.403-.414-.554-.422-.143-.008-.307-.01-.471-.01-.164 0-.43.061-.656.307-.225.246-.86.84-.86 2.05s.88 2.378 1.004 2.542c.123.164 1.733 2.646 4.198 3.71.586.253 1.044.405 1.4.518.589.187 1.125.161 1.549.098.472-.07 1.455-.594 1.66-1.168.205-.574.205-1.066.144-1.168-.062-.103-.225-.164-.471-.287z"
                  />
                </svg>
              </div>
              <div>
                <h4 class="text-sm font-bold tracking-tight">
                  WhatsApp Support Probus
                </h4>
                <p class="text-[11px] text-emerald-100/90 mt-0.5">
                  Pilih kontak di bawah untuk langsung chat:
                </p>
              </div>
            </div>

            <div class="flex items-center gap-1.5">
              <span
                class="text-[10px] bg-white/20 text-white font-bold px-2 py-0.5 rounded-full shrink-0"
              >
                3 Kontak
              </span>
              <button
                type="button"
                @click="isOpen = false"
                class="w-7 h-7 rounded-lg text-white/80 hover:text-white hover:bg-white/20 flex items-center justify-center transition-colors cursor-pointer"
                title="Tutup Menu"
              >
                <X class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>

        <!-- 3 WhatsApp Contacts List -->
        <div class="p-2 divide-y divide-slate-100 max-h-[60vh] overflow-y-auto">
          <div
            v-for="contact in contacts"
            :key="contact.id"
            @click="openWhatsApp(contact)"
            class="p-3 rounded-xl hover:bg-slate-50 cursor-pointer transition-all duration-150 flex items-start gap-3 group border border-transparent hover:border-slate-200/80 active:scale-[0.99]"
          >
            <!-- Icon Box -->
            <div
              class="w-10 h-10 rounded-xl bg-slate-100 text-slate-700 flex items-center justify-center shrink-0 mt-0.5 group-hover:bg-emerald-50 group-hover:text-emerald-600 transition-colors shadow-2xs"
            >
              <component :is="contact.icon" class="w-4.5 h-4.5 stroke-[2]" />
            </div>

            <!-- Contact Details -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between gap-1">
                <h5
                  class="text-xs sm:text-sm font-bold text-slate-900 group-hover:text-emerald-700 transition-colors truncate"
                >
                  {{ contact.name }}
                </h5>
                <span
                  :class="[
                    'text-[10px] font-bold px-2 py-0.5 rounded border shrink-0',
                    contact.badgeClass,
                  ]"
                >
                  {{ contact.badge }}
                </span>
              </div>
              <p class="text-[11px] text-slate-500 mt-0.5 line-clamp-1">
                {{ contact.role }}
              </p>
              <div
                class="mt-1.5 flex items-center gap-1.5 text-xs font-semibold text-emerald-700"
              >
                <span>{{ contact.numberDisplay }}</span>
                <ExternalLink
                  class="w-3 h-3 text-emerald-500 group-hover:translate-x-0.5 transition-transform"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- Footer Operational Hours Note -->
        <div
          class="p-3 bg-slate-50 border-t border-slate-100 text-[11px] text-slate-500 text-center flex items-center justify-center gap-1.5"
        >
          <span
            class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"
          ></span>
          <span
            >Tim support siap melayani Anda pada jam kerja operasional.</span
          >
        </div>
      </div>
    </Transition>

    <!-- Floating Action Button (FAB) at Bottom-Right -->
    <div class="relative group">
      <!-- Tooltip label (Visible on desktop hover when closed) -->
      <div
        v-if="!isOpen"
        class="absolute right-full mr-3 top-1/2 -translate-y-1/2 hidden md:flex items-center gap-1 px-3 py-1.5 rounded-xl bg-slate-900 text-white text-xs font-semibold shadow-lg whitespace-nowrap opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none"
      >
        <span>Chat WhatsApp Support</span>
        <div
          class="w-2 h-2 bg-slate-900 rotate-45 absolute -right-1 top-1/2 -translate-y-1/2"
        ></div>
      </div>

      <!-- Main Button -->
      <button
        @click.stop="toggleOpen"
        type="button"
        class="w-14 h-14 rounded-full bg-gradient-to-tr from-emerald-600 to-emerald-500 hover:from-emerald-500 hover:to-emerald-400 text-white flex items-center justify-center shadow-lg shadow-emerald-500/30 hover:shadow-xl hover:shadow-emerald-500/40 hover:scale-105 active:scale-95 transition-all duration-200 cursor-pointer focus:outline-none focus:ring-4 focus:ring-emerald-400/30"
        :aria-expanded="isOpen"
        title="WhatsApp Support"
      >
        <!-- Icon state: WhatsApp Icon or Close X -->
        <Transition
          mode="out-in"
          enter-active-class="transition duration-150 ease-out"
          enter-from-class="transform opacity-0 rotate-45 scale-75"
          enter-to-class="transform opacity-100 rotate-0 scale-100"
          leave-active-class="transition duration-100 ease-in"
          leave-from-class="transform opacity-100 rotate-0 scale-100"
          leave-to-class="transform opacity-0 -rotate-45 scale-75"
        >
          <X v-if="isOpen" class="w-6 h-6 stroke-[2.5]" />
          <svg
            v-else
            viewBox="0 0 24 24"
            class="w-7 h-7 fill-white drop-shadow-xs"
          >
            <path
              d="M12.031 2C6.511 2 2.016 6.496 2.016 12.016c0 1.988.583 3.918 1.688 5.568L2.016 22l4.57-1.636c1.604.996 3.473 1.652 5.445 1.652 5.52 0 10.016-4.496 10.016-10.016 0-5.52-4.496-10-10.016-10zm.016 18.234c-1.745 0-3.414-.492-4.852-1.422l-.348-.225-2.871 1.028.98-2.784-.243-.377c-1.043-1.522-1.6-3.328-1.6-5.454 0-4.524 3.68-8.204 8.203-8.204 4.524 0 8.204 3.68 8.204 8.204 0 4.523-3.68 8.234-8.204 8.234zm4.496-6.148c-.246-.123-1.455-.718-1.68-.8-.225-.082-.389-.123-.553.123-.164.246-.635.8-.779.963-.143.164-.287.184-.533.061-.246-.123-1.037-.383-1.977-1.221-.73-.652-1.223-1.457-1.367-1.703-.143-.246-.015-.379.108-.502.111-.111.246-.287.369-.43.123-.144.164-.246.246-.41.082-.164.041-.308-.02-.43-.062-.124-.554-1.334-.759-1.826-.2-.48-.403-.414-.554-.422-.143-.008-.307-.01-.471-.01-.164 0-.43.061-.656.307-.225.246-.86.84-.86 2.05s.88 2.378 1.004 2.542c.123.164 1.733 2.646 4.198 3.71.586.253 1.044.405 1.4.518.589.187 1.125.161 1.549.098.472-.07 1.455-.594 1.66-1.168.205-.574.205-1.066.144-1.168-.062-.103-.225-.164-.471-.287z"
            />
          </svg>
        </Transition>

        <!-- Notification Pulse Indicator when closed -->
        <span
          v-if="!isOpen"
          class="absolute -top-1 -right-1 w-4 h-4 bg-white rounded-full flex items-center justify-center shadow-xs"
        >
          <span
            class="w-2.5 h-2.5 bg-emerald-500 rounded-full animate-pulse"
          ></span>
        </span>
      </button>
    </div>
  </div>
</template>
