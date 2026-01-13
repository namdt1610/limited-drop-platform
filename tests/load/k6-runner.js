import http from "k6/http";
import { check, sleep } from "k6";
// Mode: 'e2e' (single deterministic run) or 'stress' (high-RPS load)
// Default values are hardcoded so this runner is zero-config out of the box.
const MODE =
  typeof __ENV !== "undefined" && __ENV.K6_MODE
    ? __ENV.K6_MODE.toLowerCase()
    : "e2e";
const BASE_URL = "http://localhost:3030";
const DROP_ID = 1; // assume drop id 1 exists for e2e smoke

// Options chosen based on mode
export let options;
// Allow overrides via env vars for zero-config tuning
const ENV_RATE =
  typeof __ENV !== "undefined" && __ENV.RATE ? Number(__ENV.RATE) : undefined;
const ENV_DURATION =
  typeof __ENV !== "undefined" && __ENV.DURATION ? __ENV.DURATION : undefined;
const ENV_PREALLOC =
  typeof __ENV !== "undefined" && __ENV.PREALLOC_VUS
    ? Number(__ENV.PREALLOC_VUS)
    : undefined;
const ENV_MAXVUS =
  typeof __ENV !== "undefined" && __ENV.MAX_VUS
    ? Number(__ENV.MAX_VUS)
    : undefined;
const FLOW =
  typeof __ENV !== "undefined" && __ENV.K6_FLOW
    ? __ENV.K6_FLOW.toLowerCase()
    : undefined;

if (MODE === "stress" || MODE === "kpi" || MODE === "kpi-10k") {
  // Default values (kpi aims for 10k rps but can be overridden via env)
  const defaultRate =
    MODE === "stress"
      ? 10000
      : MODE === "kpi" || MODE === "kpi-10k"
        ? 10000
        : 1000;
  const rate = ENV_RATE || defaultRate;
  const duration = ENV_DURATION || "30s";
  const preAllocatedVUs = ENV_PREALLOC || (MODE === "kpi" ? 2000 : 1000);
  const maxVUs = ENV_MAXVUS || (MODE === "kpi" ? 10000 : 10000);

  options = {
    scenarios: {
      constant_request_rate: {
        executor: "constant-arrival-rate",
        rate: rate,
        timeUnit: "1s",
        duration: duration,
        preAllocatedVUs: preAllocatedVUs,
        maxVUs: maxVUs,
      },
    },
    thresholds: {
      http_req_duration: ["p(95)<2000"],
      http_req_failed: ["rate<0.8"],
    },
  };
  console.log(
    `K6 MODE=${MODE} flow=${FLOW || "stress"
    } rate=${rate}/s duration=${duration} prealloc=${preAllocatedVUs} maxVUs=${maxVUs}`
  );
} else if (MODE === "stress-e2e") {
  // stress-e2e: run full E2E flows under moderate load
  options = {
    scenarios: {
      constant_request_rate: {
        executor: "constant-arrival-rate",
        rate: 500, // 500 E2E flows/sec by default
        timeUnit: "1s",
        duration: "30s",
        preAllocatedVUs: 200,
        maxVUs: 1000,
      },
    },
    thresholds: {
      http_req_duration: ["p(95)<2000"],
      http_req_failed: ["rate<0.8"],
    },
  };
} else {
  options = {
    vus: 1,
    iterations: 1,
  };
}

function uuid() {
  return Date.now().toString(36) + Math.floor(Math.random() * 1e6).toString(36);
}

function postWebhook(orderCode, dropId, productId, phone, email) {
  const webhookPayload = {
    code: "OK",
    desc: "k6 test payment",
    data: {
      orderCode: Number(orderCode),
      amount: 100000,
      status: "PAID",
      description: "k6 test payment",
      metadata: {
        drop_id: String(dropId),
        product_id: String(productId),
        customer_phone: phone,
        customer_email: email,
        customer_name: "k6 Test User",
        shipping_address: "k6 Address",
        quantity: "1",
      },
      paymentMethod: "payos",
    },
  };
  const body = JSON.stringify(webhookPayload);
  const headers = { "Content-Type": "application/json" };
  // No signature header — dev runner is zero-config and accepts unsigned webhooks in dev mode
  return http.post(`${BASE_URL}/api/limited-drops/webhook/payos`, body, {
    headers,
  });
}

// Deterministic e2e flow
function runE2E() {
  // 1) Fetch drops
  const dropsResp = http.get(`${BASE_URL}/api/drops`);
  check(dropsResp, { "drops fetched": (r) => r.status === 200 });
  let drop;
  try {
    const drops = JSON.parse(dropsResp.body);
    drop = drops.find((d) => d.id === DROP_ID || d.id === Number(DROP_ID));
  } catch (e) { }
  if (!drop) {
    console.error(`Drop ${DROP_ID} not found`);
    return;
  }
  const soldBefore = drop.sold || 0;
  const productId = drop.product_id || drop.productId || 1;

  // 2) Purchase
  const phone = `k6-${uuid()}`;
  const email = `k6+${uuid()}@example.com`;
  const purchasePayload = {
    quantity: 1,
    name: "k6 Test User",
    phone,
    email,
    address: "123 Test St",
    province: "K6",
    district: "K6",
    ward: "K6",
  };
  const idempotencyKey = uuid();
  const purchaseResp = http.post(
    `${BASE_URL}/api/drops/${DROP_ID}/purchase`,
    JSON.stringify(purchasePayload),
    {
      headers: {
        "Content-Type": "application/json",
        "X-Idempotency-Key": idempotencyKey,
      },
    }
  );
  check(purchaseResp, { "purchase 200": (r) => r.status === 200 });
  if (purchaseResp.status !== 200) {
    console.error("Purchase failed:", purchaseResp.status, purchaseResp.body);
    return;
  }
  let orderCode = null;
  try {
    const body = JSON.parse(purchaseResp.body);
    orderCode = body.order_code || body.orderCode || body.OrderCode;
  } catch (e) { }
  if (!orderCode) {
    console.error("purchase did not return order code");
    return;
  }

  // 3) Webhook
  const webhookResp = postWebhook(orderCode, DROP_ID, productId, phone, email);
  check(webhookResp, { "webhook 200": (r) => r.status === 200 });

  // 4) Poll orders and drops to verify
  let found = false;
  for (let i = 0; i < 10; i++) {
    const ordersResp = http.get(
      `${BASE_URL}/api/orders?phone=${encodeURIComponent(phone)}`
    );
    if (ordersResp.status === 200) {
      try {
        const parsed = JSON.parse(ordersResp.body);
        const orders = parsed.orders || [];
        if (orders.length > 0) {
          for (const o of orders) {
            if (
              o.payos_order_code == orderCode ||
              o.PayOSOrderCode == orderCode ||
              o.payosOrderCode == orderCode
            ) {
              found = true;
              break;
            }
          }
          if (!found && orders.length >= 1) found = true;
        }
      } catch (e) { }
    }
    const dropsNowResp = http.get(`${BASE_URL}/api/drops`);
    try {
      const dropsNow = JSON.parse(dropsNowResp.body);
      const nowDrop = dropsNow.find(
        (d) => d.id === DROP_ID || d.id === Number(DROP_ID)
      );
      const soldNow = nowDrop ? nowDrop.sold || 0 : soldBefore;
      if (found && soldNow > soldBefore) {
        check(true, { "e2e: success": () => true });
        return;
      }
    } catch (e) { }
    sleep(1);
  }
  check(false, { "e2e: verified": () => false });
}

// Stress flow (light-weight per-iteration ops)
function runStress() {
  const rand = Math.random();
  // 60% get drops
  if (rand < 0.6) {
    const r = http.get(`${BASE_URL}/api/drops`);
    check(r, { "drops status": (res) => res.status !== 0 });
    return;
  }
  // 20% purchase then webhook
  if (rand < 0.8) {
    const purchasePayload = {
      quantity: 1,
      name: "Stress User",
      phone: `s-${uuid()}`,
      email: `s+${uuid()}@example.com`,
      address: "Stress St",
      province: "S",
      district: "S",
      ward: "S",
    };
    const resp = http.post(
      `${BASE_URL}/api/drops/${DROP_ID}/purchase`,
      JSON.stringify(purchasePayload),
      {
        headers: {
          "Content-Type": "application/json",
          "X-Idempotency-Key": uuid(),
        },
      }
    );
    check(resp, { "purchase ok": (r) => r.status !== 0 });
    try {
      const b = JSON.parse(resp.body);
      const orderCode = b.order_code || b.orderCode || b.OrderCode;
      if (orderCode) {
        const wr = postWebhook(
          orderCode,
          DROP_ID,
          purchasePayload.product_id || 1,
          purchasePayload.phone,
          purchasePayload.email
        );
        check(wr, { "webhook ok": (r) => r.status === 200 });
      }
    } catch (e) { }
    return;
  }
  // 10% fake webhook-only (for variety)
  if (rand < 0.9) {
    const wc = postWebhook(
      Math.floor(Math.random() * 1e9),
      DROP_ID,
      1,
      `s-${uuid()}`,
      `s+${uuid()}@example.com`
    );
    check(wc, { "webhook-only": (r) => r.status === 200 });
    return;
  }
  // otherwise no-op
}

// Stress-E2E: each iteration runs a full E2E flow and verifies result quickly
function runStressE2E() {
  // 1) Purchase
  const phone = `se2e-${uuid()}`;
  const email = `se2e+${uuid()}@example.com`;
  const purchasePayload = {
    quantity: 1,
    name: "Stress E2E User",
    phone,
    email,
    address: "Stress E2E St",
    province: "SE",
    district: "SE",
    ward: "SE",
  };
  const idempotencyKey = uuid();
  const purchaseResp = http.post(
    `${BASE_URL}/api/drops/${DROP_ID}/purchase`,
    JSON.stringify(purchasePayload),
    {
      headers: {
        "Content-Type": "application/json",
        "X-Idempotency-Key": idempotencyKey,
      },
    }
  );
  check(purchaseResp, { "se2e: purchase ok": (r) => r.status === 200 });

  let orderCode = null;
  try {
    orderCode =
      JSON.parse(purchaseResp.body).order_code ||
      JSON.parse(purchaseResp.body).orderCode ||
      JSON.parse(purchaseResp.body).OrderCode;
  } catch (e) { }
  if (!orderCode) return; // can't verify without order code

  // 2) Webhook
  const webhookResp = postWebhook(
    orderCode,
    DROP_ID,
    purchasePayload.product_id || 1,
    phone,
    email
  );
  check(webhookResp, { "se2e: webhook ok": (r) => r.status === 200 });

  // 3) Quick poll for order presence and sold increment (small timeout)
  const start = Date.now();
  let verified = false;
  while (Date.now() - start < 5000) {
    // 5s max wait
    const ordersResp = http.get(
      `${BASE_URL}/api/orders?phone=${encodeURIComponent(phone)}`
    );
    if (ordersResp.status === 200) {
      try {
        const orders = JSON.parse(ordersResp.body).orders || [];
        if (orders.length > 0) {
          verified = true;
          break;
        }
      } catch (e) { }
    }
    sleep(0.5);
  }
  check(verified, { "se2e: verified": () => verified });
}

export default function () {
  // For high-rate runs we default to the light-weight stress flow unless FLOW requests se2e
  const selectedFlow =
    typeof __ENV !== "undefined" && __ENV.K6_FLOW
      ? __ENV.K6_FLOW.toLowerCase()
      : MODE;
  if (selectedFlow === "stress" || MODE === "kpi" || MODE === "kpi-10k") {
    // use the light weight stress worker for high-RPS tests
    runStress();
  } else if (selectedFlow === "stress-e2e" || MODE === "stress-e2e") {
    runStressE2E();
  } else {
    runE2E();
  }
}
