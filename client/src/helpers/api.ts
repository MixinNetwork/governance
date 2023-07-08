import axios, { type AxiosResponse } from 'axios';
import axiosRetry from 'axios-retry';
import { ResponseError } from '@mixin.dev/mixin-node-sdk';
import type { NodeResponse } from '@/types';
import { API_URL } from './constant';

axios.defaults.headers.post['Content-Type'] = 'application/json';

export const initSafeClient = () => {
  const ins = axios.create({
    baseURL: API_URL,
    timeout: 5000,
  });

  ins.interceptors.response.use(async (res: AxiosResponse) => {
    const { data, error } = res.data;
    if (error)
      throw new ResponseError(error.code, error.description, error.status, error.extra, '', error);
    return data;
  });

  axiosRetry(ins, {
    retries: 5,
    shouldResetTimeout: true,
    retryDelay: () => 500,
  });

  return {
    listNodes: (): Promise<NodeResponse[]> => ins.get('/nodes'),

    register: (extra: string): Promise<NodeResponse> =>
      ins.post(`/nodes`, {
        extra,
      }),
  };
};
