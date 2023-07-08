<script setup lang="ts">
import { onMounted } from 'vue';
import { storeToRefs } from 'pinia';
import { NCard, NSkeleton, NCollapse, NCollapseItem, NConfigProvider } from 'naive-ui';
import { useNodeStore } from '@/stores/node';
import { NAIVE_THEMES, API_URL, BOT_PUBLIC_KEY } from '@/helpers/constant';

const nodeStore = useNodeStore();
const { loading, nodes } = storeToRefs(nodeStore);
const { fetchNodes } = nodeStore;

onMounted(async () => {
  await fetchNodes();
});

const steps = [
  {
    title: 'Generate Custodian Address and Private Key',
    text: `
      git clone <a style="color: blue;" target="_blank" href="https://github.com/MixinNetwork/mixin.git">https://github.com/MixinNetwork/mixin.git</a><br/>
      cd mixin<br/>
      go build<br/><br/>

      ./mixin createaddress -public<br/>
      # address: XINRXkrW1...tX6d<br/>
      # private view key: 708a8db3...0b2100<br/>
      # private spend key: c0619ce9...54be0c<br/>
    `,
  },
  {
    title: 'Register Custodian Node',
    text: 'Submit new generated custodian node address with the information of your Mixin node, then pay 100XIN fee. You will see your node in active nodes list if success.',
  },
  {
    title: 'Transfer App',
    text: `
      git clone <a style="color: blue;" target="_blank" href="https://github.com/MixinNetwork/governance.git">https://github.com/MixinNetwork/governance.git</a><br/>
      cd governance<br/>
      go build<br/><br/>

      # You can find the keystore and publickey of your custodian node in the active nodes list after registration.<br/>
      governance migrate -k keystore -s custodianspendkey -p publickey -u receiver<br/>
        -k Encrypted bot keystore, base64<br/>
        -p Public key used to decrypt keystore, ${BOT_PUBLIC_KEY} <br/>
        -s Custodian private spend key, generated in step 1<br/>
        -u Uuid of the Mixin Messenger user who will receive the app<br/>
    `,
  },
  {
    title: 'Prepare Safe Node Config',
    text: `Download config template from <a style="color: blue;" target="_blank" href="${API_URL}/template">${API_URL}/template</a> and modify it with your app info.`,
  },
];
</script>

<template>
  <main class="px-10 py-20">
    <RouterLink
      :to="`/register`"
      class="block mb-10 py-2 px-3 w-fit bg-[#4B7CDD] text-white rounded"
    >
      Register Custodian Node
    </RouterLink>

    <n-config-provider :theme-overrides="NAIVE_THEMES" class="flex justify-between">
      <n-card
        class="w-[49%]"
        title="Active Nodes"
        size="huge"
        :segmented="{
          content: true,
        }"
      >
        <template v-if="loading">
          <div v-for="i of 5" :key="i" :class="['py-4', i !== 5 ? 'border-b-[1px]' : '']">
            <n-skeleton text class="h-[22.4px]" />
          </div>
        </template>
        <div v-else>
          <n-collapse
            v-if="nodes.length > 0"
            v-for="(n, i) of nodes"
            :key="i"
            :default-expanded-names="n.kernel_id"
            accordion
          >
            <n-collapse-item
              :title="n.kernel_id"
              :name="i"
              :class="['py-4', i !== nodes.length - 1 ? 'border-b-[1px]' : '']"
            >
              <div class="px-5">
                <div>
                  <h3 class="font-bold text-base">Custodian</h3>
                  <div>{{ n.custodian }}</div>
                </div>
                <div>
                  <h3 class="font-bold text-base">App ID</h3>
                  <div>{{ n.app_id }}</div>
                </div>
                <div>
                  <h3 class="font-bold text-base">Keystore</h3>
                  <div>{{ n.keystore }}</div>
                </div>
              </div>
            </n-collapse-item>
          </n-collapse>
          <div v-else>No Active Nodes Yet</div>
        </div>
      </n-card>

      <n-card
        class="w-[49%]"
        title="Register Steps"
        size="huge"
        :segmented="{
          content: true,
        }"
      >
        <div v-for="(s, i) of steps" :key="i" class="mb-5">
          <h3 class="mb-1 font-bold text-base">{{ `${i + 1}. ${s.title}` }}</h3>
          <div v-html="s.text" class="text-base"></div>
        </div>
      </n-card>
    </n-config-provider>
  </main>
</template>
