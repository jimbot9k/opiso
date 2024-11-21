<template>
  <div v-if="visible" class="fixed bottom-4 right-4 bg-green-800 text-white px-4 py-2 rounded shadow-lg">
    {{ message }}
  </div>
</template>

<script lang="ts" setup>
import { onUnmounted, ref, watch } from 'vue';

const props = defineProps({
  message: {
    type: String,
    required: true,
  },
  duration: {
    type: Number,
    default: 3000,
  },
});

const visible = ref(false);
let timeoutId: number | null = null;

watch(() => props.message, (newMessage) => {
  if (!newMessage) {
    return;
  }

  if (timeoutId !== null) {
    clearTimeout(timeoutId);
    timeoutId = null;
    visible.value = false;
  }

  visible.value = true;
  timeoutId = setTimeout(() => {
    timeoutId = null;
    visible.value = false;
  }, props.duration);
});

onUnmounted(() => {
  if (timeoutId === null) {
    return;
  }
  clearTimeout(timeoutId);
  timeoutId = null;
});
</script>