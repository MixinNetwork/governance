import { ref } from 'vue';
import { defineStore } from 'pinia';
import type { NodeResponse } from '@/types';
import { initSafeClient } from '@/helpers/api';

export const useNodeStore = defineStore('node', () => {
  const loading = ref(false);
  const nodes = ref<NodeResponse[]>([]);

  const client = initSafeClient();

  const fetchNodes = async () => {
    loading.value = true;
    nodes.value = await client.listNodes();
    loading.value = false;
  };

  return {
    loading,
    nodes,
    fetchNodes,
  };
});
