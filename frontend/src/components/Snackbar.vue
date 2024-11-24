<template>
  <div v-if="visible && message" class="fixed bottom-4 right-4 bg-green-800 text-white px-4 py-2 rounded shadow-lg">
    {{ message }}
  </div>
</template>

<script lang="ts" setup>
import { onUnmounted, ref, watch } from 'vue';

const props = defineProps({
  message: {
    type: [String, null],
    required: true,
  },
  duration: {
    type: Number,
    default: 3000,
  },
});

const visible = ref(false);
const message = ref<string | null>(null);

let timeoutId: number | undefined;

watch(() => props.message, (newMessage) => {
  if (!newMessage) {
    message.value = null;
    visible.value = false;
    snackbarClearTimeoutIfExists();
    return;
  }

  snackbarClearTimeoutIfExists();
  visible.value = true;
  message.value = props.message;
  timeoutId = setTimeout(() => {
    visible.value = false;
    message.value = null;
  }, props.duration);
});

onUnmounted(() => {
  snackbarClearTimeoutIfExists()
});

const snackbarClearTimeoutIfExists = () => {
  if (timeoutId !== undefined) {
    timeoutId = undefined;
    clearTimeout(timeoutId);
  }
}
</script>