<script setup lang="ts">
import { provide, reactive } from 'vue';
import MessageForm from './components/MessageForm.vue';
import Snackbar from './components/Snackbar.vue';

const snackbarStore = reactive({
  message: '',
  duration: 3000,
  showSnackbar(message: string, duration: number = 3000) {
    this.message = message;
    this.duration = duration;
  },
});

provide('snackbar', snackbarStore);
</script>

<template>
  <header>
  </header>

  <main>
    <MessageForm></MessageForm>
    <Snackbar :message="snackbarStore.message" :duration="snackbarStore.duration"/>
  </main>
</template>

<script lang="ts" setup>

</script>
