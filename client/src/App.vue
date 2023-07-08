<script setup lang="ts">
import { RouterView } from 'vue-router';
import { NMessageProvider } from 'naive-ui';
</script>

<template>
  <n-message-provider>
    <RouterView />
  </n-message-provider>
</template>
