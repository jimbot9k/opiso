<script setup lang="ts">
import { provide, reactive } from 'vue';
import MessageForm from './components/MessageForm.vue';
import Snackbar from './components/Snackbar.vue';

const snackbarStore = reactive({
  message: <string | null>null,
  duration: 3000,
  showSnackbar(message: string | null, duration: number = 3000) {
    this.message = null;
    setTimeout(() => {
      this.message = message;
      this.duration = duration;
    }, 0) // This is in a timeout to force reactivity on duplicate messages
  },
});

provide('snackbar', snackbarStore);
</script>

<template>
  <header>
  </header>

  <main>
    <MessageForm></MessageForm>
    <Snackbar :message="snackbarStore.message" :duration="snackbarStore.duration" />
  </main>
</template>

<script lang="ts" setup>

</script>
