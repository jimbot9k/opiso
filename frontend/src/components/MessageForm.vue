<template>
  <form @submit.prevent="submitForm">
    <div>
      <textarea id="message" v-model="formData.rawMessages"
        placeholder="Write your messages here. Each word in the input counts as a message. Each message will be reversed upon submission."
        rows="5" cols="30"></textarea>
    </div>
    <button type="submit">Submit</button>
  </form>
  <ol>
    <li v-for="response in messagesResponses" :key="response">
      {{ response.reversed }}
    </li>
  </ol>
</template>

<script lang="ts">
import { defineComponent, reactive, ref } from "vue";

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
    const formData = reactive<FormData>({
      rawMessages: "",
    });
    const messagesResponses = ref<MessagesResponse[]>([])
    const loading = ref(false);
    const error = ref<string | null>(null);

    const submitForm = async () => {
      loading.value = true;
      error.value = null;

      const messages = formData.rawMessages.split(" ");
      const request: MessagesRequest = { "messages": messages };

      try {
        const response = await fetch("http://localhost:8080/reverse", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(request),
          mode: "cors",
        });

        if (!response.ok) {
          throw new Error(`Error: ${await response.json() as ErrorResponse}`);
        }
        const newResponse = await response.json() as MessagesResponse;
        messagesResponses.value.push(newResponse);
        console.log(messagesResponses.value);
        formData.rawMessages = "";
      } catch (err: any) {
        error.value = err.message;
      } finally {
        loading.value = false;
      }
    };

    return {
      formData,
      loading,
      error,
      submitForm,
      messagesResponses,
    };
  },
});
</script>
