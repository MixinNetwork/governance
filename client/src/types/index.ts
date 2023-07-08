export interface NodeResponse {
  kernel_id: string;
  custodian: string;
  payee: string;
  app_id: string;
  keystore: string;
  public_key: string;
  mixin_hash: string;
  created_at: string;
  updated_at: string;
}

export interface NodeRegisterResponse {
  hash: string;
}
