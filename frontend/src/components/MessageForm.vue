<template>
  <div class="scroll-none min-h-screen flex flex-col items-center bg-blue-100">

    <form @submit.prevent="submitForm" @keydown.ctrl.enter="submitForm" @keydown.shift.enter="submitForm"
      class="sticky top-0 flex flex-row items-end w-full gap-4 py-8 justify-center bg-blue-100 shadow-lg">
      <textarea id="message" v-model="formData.rawMessages"
        placeholder="Write your messages here. Each word in the input counts as a message. Each message will be reversed upon submission."
        rows="5" cols="60" autofocus
        class="block p-2.5 resize-none text-sm bg-white border first:border-black text-black rounded w-8/12">
      </textarea>
      <button class="text-white font-bold py-2 px-4 rounded h-12 w-60 shadow-lg"
        :class="allowSubmit() ? 'bg-gray-400' : 'bg-blue-500 hover:bg-blue-700 '" :disabled="allowSubmit()">{{ loading ?
          'Loading.' : 'Send.' }}
      </button>
    </form>

    <ol class="m-4 flex flex-col-reverse items-center gap-4">
      <li v-for="(response, index) in messagesResponses" v-on:click="copyMessages(response)"
        class="flex flex-row gap-4 text-base rounded border border-black hover:bg-blue-200 hover:cursor-copy p-4 w-fit flex-wrap shadow-lg bg-white">
        <DeleteIcon v-on:click.stop="deleteMessage(index)" class="fill-red-600"></DeleteIcon>
        <span class="break-all" v-for="word in response.reversed">{{ word }}</span>
      </li>
    </ol>
  </div>

</template>

<script lang="ts">
import { API_URL } from "@/config";
import { defineComponent, inject, reactive, ref } from "vue";
import DeleteIcon from "./DeleteIcon.vue";


interface FormData {
  rawMessages: string;
}

interface MessagesRequest {
  messages: string[];
}

interface MessagesResponse {
  reversed: string[];
}

interface ErrorResponse {
  message: string[];
}

export default defineComponent({
  setup() {
    const snackbar = inject<{ message: string; showSnackbar: (msg: string, duration?: number) => void }>('snackbar');

    const formData = reactive<FormData>({
      rawMessages: "",
    });
    const messagesResponses = ref<MessagesResponse[]>([])
    const loading = ref(false);

    const allowSubmit = (): boolean => {
      return loading.value || formData.rawMessages.length === 0
    };

    const deleteMessage = (index: number): void => {
      messagesResponses.value.splice(index, 1)
    }

    const submitForm = async () => {
      const messages = formData.rawMessages.split(" ");
      const request: MessagesRequest = { "messages": messages };
      loading.value = true;

      try {
        const requestUrl = `${API_URL}/reverse`;

        const response = await fetch(requestUrl, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(request),
          mode: "cors",
        });

        if (!response.ok) {
          throw new Error(`Error: ${(await response.json() as ErrorResponse).message}`);
        }
        const newResponse = await response.json() as MessagesResponse;
        messagesResponses.value.push(newResponse);
        formData.rawMessages = "";
        window.scrollTo(0, 0);

        const entries = performance.getEntriesByName(requestUrl);
        performance.clearResourceTimings();
        if (entries.length > 0) {
          const timing = entries[0];
          const duration = timing.duration
          snackbar?.showSnackbar(`Messaged reversed successfully in ${duration.toFixed(2)} ms`, 3000);
        }
      } catch (err: any) {
        snackbar?.showSnackbar(err.message, 3000);
      } finally {
        loading.value = false;
      }
    };

    const copyMessages = (messages: MessagesResponse) => {
      navigator.clipboard.writeText(messages.reversed.join(" "));
      snackbar?.showSnackbar(`Message copied`, 3000);
    }

    return {
      formData,
      loading,
      messagesResponses,
      submitForm,
      copyMessages,
      allowSubmit,
      deleteMessage
    };
  },
  components: {
    DeleteIcon
  }
});
</script>
