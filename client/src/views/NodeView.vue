<script setup lang="ts">
import { computed, reactive, ref } from 'vue';
import { v4 } from 'uuid';
import { NCard, useMessage } from 'naive-ui';
import Spinner from '@/components/Common/Spinner.vue';
import { buildExtra } from '@/helpers/register';
import { initSafeClient } from '@/helpers/api';
import { BOT_ID, FEE_ASSET_ID } from '@/helpers/constant';

const message = useMessage();

const loading = ref(false);

const state = reactive({
  node_id: '',
  custodian: '',
  payee: '',
  signerSpendKey: '',
  payeeSpendKey: '',
  custodianSpendKey: '',
});
const isNoEmpty = computed(
  () =>
    !!state.node_id &&
    !!state.custodian &&
    !!state.payee &&
    !!state.payeeSpendKey &&
    !!state.signerSpendKey &&
    !!state.custodianSpendKey,
);

const useRegister = async () => {
  if (!isNoEmpty.value) {
    return;
  }
  loading.value = true;

  const client = initSafeClient();
  try {
    const extra = await buildExtra(state);
    console.log(extra);
    const node = await client.register(extra);

    const id = v4();
    location.href = `http://mixin.one/pay?recipient=${BOT_ID}&asset=${FEE_ASSET_ID}&trace=${id}&amount=100&memo=${node.mixin_hash}`;
  } catch (e: any) {
    let msg: string;
    if (e.description) msg = e.description;
    else msg = e.message;

    message.error(msg, {
      closable: true,
      duration: 5000,
    });
  }
};
</script>

<template>
  <main class="py-20 mx-auto w-2/3">
    <n-card
      title="Register Custodian Node"
      size="huge"
      :segmented="{
        content: true,
      }"
    >
      <div class="text-base">
        <div class="flex justify-between items-center h-20">
          <label class="w-1/6" for="node_id">Node ID</label>
          <input class="p-3 w-4/6 h-12" type="text" id="node_id" v-model="state.node_id" />
        </div>
        <div class="flex justify-between items-center h-20">
          <label class="w-1/6" for="custodian">Custodian</label>
          <input class="p-3 w-4/6 h-12" type="text" id="custodian" v-model="state.custodian" />
        </div>
        <div class="flex justify-between items-center h-20">
          <label class="w-1/6" for="payee">Payee</label>
          <input class="p-3 w-4/6 h-12" type="text" id="payee" v-model="state.payee" />
        </div>
        <div class="flex justify-between items-center h-20">
          <label class="w-1/6" for="signerSpendKey">Signer Spend Key</label>
          <input
            class="p-3 w-4/6 h-12"
            type="text"
            id="signerSpendKey"
            v-model="state.signerSpendKey"
          />
        </div>
        <div class="flex justify-between items-center h-20">
          <label class="w-1/6" for="payeeSpendKey">Payee Spend Key</label>
          <input
            class="p-3 w-4/6 h-12"
            type="text"
            id="payeeSpendKey"
            v-model="state.payeeSpendKey"
          />
        </div>
        <div class="flex justify-between items-center h-20">
          <label class="w-1/6" for="custodianSpendKey">Custodian Spend Key</label>
          <input
            class="p-3 w-4/6 h-12"
            type="text"
            id="custodianSpendKey"
            v-model="state.custodianSpendKey"
          />
        </div>
      </div>
      <template #action>
        <div class="flex justify-center">
          <button
            :class="[
              'flex justify-center py-2 min-w-[86px] text-white rounded',
              isNoEmpty ? 'bg-primary' : 'bg-black/[.4] cursor-not-allowed',
            ]"
            @click="useRegister"
            :disabled="!isNoEmpty"
          >
            <Spinner v-if="loading" />
            <template v-else>Register</template>
          </button>
        </div>
      </template>
    </n-card>
  </main>
</template>
