import { base64RawURLEncode } from '@mixin.dev/mixin-node-sdk';
import { SHA3 } from 'sha3';
import * as ed from '@noble/ed25519';
import { utils } from 'ethers';

const MainNetworkId = 'XIN';

type ExtK = { prefix: Buffer; scalar: bigint; pointBytes: Buffer };

const n2b_32LE = (num: bigint) =>
  ed.etc.hexToBytes(num.toString(16).padStart(32 * 2, '0')).reverse();
const b2n_LE = (b: Buffer): bigint => BigInt('0x' + ed.etc.bytesToHex(b.reverse()));
const modL_LE = (hash: Buffer): bigint => ed.etc.mod(b2n_LE(hash), ed.CURVE.n);

const sha512Buffer = async (...messages: Uint8Array[]) => {
  const u8 = await ed.etc.sha512Async(...messages);
  return Buffer.from(u8);
};

const getExtendedPublicKey = async (priv: Buffer) => {
  const hashed = await sha512Buffer(priv);
  const prefix = hashed.subarray(32, 64);

  const scalar = modL_LE(priv); // modular division over curve order
  const point = ed.ExtendedPoint.BASE.mul(scalar); // public key point
  const pointBytes = Buffer.from(point.toRawBytes()); // point serialized to Uint8Array
  return { prefix, scalar, pointBytes };
};

const _sign = (e: ExtK, rBytes: Buffer, msg: Buffer) => {
  const { pointBytes: P, scalar: s } = e;
  const r = modL_LE(rBytes); // r was created outside, reduce it modulo L
  const R = ed.ExtendedPoint.BASE.mul(r).toRawBytes(); // R = [r]B
  const hashable = Buffer.concat([R, P, msg]);
  const finish = (hashed: Buffer) => {
    const S = ed.etc.mod(r + modL_LE(hashed) * s, ed.CURVE.n); // S = (r + k * s) mod L; 0 <= s < l
    return Buffer.concat([R, n2b_32LE(S)]);
  };
  return { hashable, finish };
};

const getAddressBuffer = (addr: string) => {
  if (!addr.startsWith(MainNetworkId)) throw new Error('invalid address network');

  const data = Buffer.from(utils.base58.decode(addr.slice(MainNetworkId.length)));
  if (data.byteLength !== 68) throw new Error('invalid address format');

  const msg = Buffer.concat([Buffer.from(MainNetworkId), data.subarray(0, 64)]);

  const checksum = new SHA3(256).update(msg).digest('hex');
  if (!Buffer.from(checksum, 'hex').subarray(0, 4).equals(data.subarray(64)))
    throw new Error('invalid address checksum');

  return data.subarray(0, 64);
};
const getNodeBuffer = (node: string) => {
  const hash = Buffer.from(node, 'hex');
  if (hash.byteLength !== 32) throw new Error(`invalid node length ${hash.byteLength}`);
  return hash;
};

const sign = async (msg: Buffer, priv: Buffer) => {
  const e = await getExtendedPublicKey(priv);
  const rBytes = await sha512Buffer(e.prefix, msg);
  const { hashable, finish } = _sign(e, rBytes, msg);
  const hash = await sha512Buffer(hashable);
  return finish(hash);
};

export const buildExtra = async (data: {
  node_id: string;
  custodian: string;
  payee: string;
  signerSpendKey: string;
  payeeSpendKey: string;
  custodianSpendKey: string;
}) => {
  const prefix = Buffer.from([1]);
  try {
    const msg = Buffer.concat([
      prefix,
      getAddressBuffer(data.custodian),
      getAddressBuffer(data.payee),
      getNodeBuffer(data.node_id),
    ]);

    const signerSpendKeyBuffer = Buffer.from(data.signerSpendKey, 'hex');
    const signerSignature = await sign(msg, signerSpendKeyBuffer);

    const payeeSpendKeyBuffer = Buffer.from(data.payeeSpendKey, 'hex');
    const payeeSignature = await sign(msg, payeeSpendKeyBuffer);

    const custodianSpendKeyBuffer = Buffer.from(data.custodianSpendKey, 'hex');
    const custodianSignature = await sign(msg, custodianSpendKeyBuffer);

    return base64RawURLEncode(
      Buffer.concat([msg, signerSignature, payeeSignature, custodianSignature]),
    );
  } catch (e: any) {
    console.log(e);
    throw new Error(e);
  }
};
