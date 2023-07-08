CREATE TABLE IF NOT EXISTS nodes (
  custodian   VARCHAR NOT NULL,
  payee       VARCHAR NOT NULL,
  kernel_id   VARCHAR NOT NULL,
  app_id      VARCHAR,
  mixin_hash  VARCHAR,
  keystore    VARCHAR NOT NULL,
  public_key  VARCHAR NOT NULL,
  created_at  TIMESTAMP NOT NULL,
  updated_at  TIMESTAMP NOT NULL,
  PRIMARY KEY ('custodian')
);

CREATE UNIQUE INDEX IF NOT EXISTS nodes_by_payee ON nodes(payee);
CREATE UNIQUE INDEX IF NOT EXISTS nodes_by_kernel_id ON nodes(kernel_id);
CREATE UNIQUE INDEX IF NOT EXISTS nodes_by_app_id ON nodes(app_id);
CREATE UNIQUE INDEX IF NOT EXISTS nodes_by_hash ON nodes(mixin_hash);
