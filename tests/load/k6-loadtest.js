import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter } from 'k6/metrics';

// Custom metrics
const successfulPurchases = new Counter('successful_purchases');
const soldOutErrors = new Counter('sold_out_errors');

export const options = {
  stages: [
    { duration: '10s', target: 500 },  // Warm up
    { duration: '30s', target: 2000 }, // RAMP TO 2000 - BREAKING POINT
    { duration: '10s', target: 0 },    // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'], // Allow higher latency under extreme stress
    // Removed http_req_failed threshold as 400s are expected
  },
};

const BASE_URL = 'http://localhost:3030';
const DROP_ID = 1;

export default function () {
  // CORE FLOW STRESS TEST (100% Writes)

  // WRITE WORKLOAD
  const purchasePayload = JSON.stringify({
    quantity: 1,
    name: `K6 User ${__VU}`,
    phone: '0987654321',
    email: `k6user${__VU}@test.com`,
    address: '123 K6 St',
    province: 'Hanoi',
    district: 'Cau Giay',
    ward: 'Dich Vong',
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  // 1. Attempt Purchase
  const res = http.post(`${BASE_URL}/api/drops/${DROP_ID}/purchase`, purchasePayload, params);

  // Check if purchase request was successful (200 OK)
  if (res.status === 200) {
    try {
      const body = JSON.parse(res.body);
      const orderCode = body.order_code || body.orderCode;

      if (orderCode) {
        // 2. Simulate Webhook (Payment Success)
        const webhookPayload = JSON.stringify({
          code: '00',
          data: {
            orderCode: orderCode,
            amount: 10000,
            status: 'PAID',
          },
        });

        const wbRes = http.post(`${BASE_URL}/api/limited-drops/webhook/payos`, webhookPayload, params);

        if (wbRes.status === 200) {
          successfulPurchases.add(1);
        } else {
          soldOutErrors.add(1);
        }
      }
    } catch (e) {
      console.error('Failed to parse response:', e);
    }
  } else if (res.status === 400) {
    // 400 Bad Request usually means Sold Out - THIS IS EXPECTED IN STRESS TEST
    soldOutErrors.add(1);
  } else {
    // 5xx errors are REAL failures
    console.error(`Unexpected status: ${res.status}`);
  }

  // sleep(0.1); // Removed sleep to test max throughput
}
