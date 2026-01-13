// API utilities for Alpine.js version
const API_BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:3030";

export async function fetchApi(endpoint, options = {}) {
  const { data, ...fetchOptions } = options;

  // Prepare headers
  const headers = {
    "Content-Type": "application/json",
    ...(fetchOptions.headers || {}),
  };

  // Prepare request options
  const requestOptions = {
    ...fetchOptions,
    headers,
  };

  // Handle request body
  if (data) {
    requestOptions.body = JSON.stringify(data);
  }

  // Make the request
  const url = `${API_BASE_URL}${endpoint}`;
  const response = await fetch(url, requestOptions);

  // Handle response
  if (!response.ok) {
    let errorData;
    try {
      errorData = await response.json();
    } catch {
      errorData = { message: await response.text() };
    }
    throw new Error(
      errorData.message || errorData.error || `API Error: ${response.status}`
    );
  }

  // Try to parse JSON response
  try {
    return await response.json();
  } catch {
    // If response is not JSON, return empty object for success responses
    return {};
  }
}

// Drop APIs
export async function getDrops() {
  const response = await fetchApi("/api/drops");
  return Array.isArray(response) ? response : response.drops || [];
}

export async function getDropStatus(dropId) {
  const response = await fetchApi(`/api/drops/${dropId}/status`);
  return response;
}

export async function purchaseDrop(dropId, purchaseData) {
  const response = await fetchApi(`/api/drops/${dropId}/purchase`, {
    method: "POST",
    data: purchaseData,
  });
  return response;
}

// Payment APIs
export async function trackOrder(trackData) {
  const response = await fetchApi("/api/payment/track-order", {
    method: "POST",
    data: trackData,
  });
  return response.data;
}

export async function verifyPayment(orderCode) {
  const response = await fetchApi(`/api/payment/payos/verify/${orderCode}`);
  return response.data;
}

export async function cancelPayment(orderCode) {
  const response = await fetchApi(`/api/payment/payos/cancel/${orderCode}`, {
    method: "POST",
  });
  return response;
}

// Symbicode verification
export async function verifySymbicode(code) {
  const response = await fetchApi(`/api/symbicode/verify`, {
    method: "POST",
    data: { code },
  });
  return response;
}

// Products
export async function getProducts(params = {}) {
  const searchParams = new URLSearchParams();
  if (params.page) searchParams.set("page", String(params.page));
  if (params.limit) searchParams.set("limit", String(params.limit));
  if (params.search) searchParams.set("search", params.search);
  if (params.sort) searchParams.set("sort", params.sort);
  if (params.minPrice) searchParams.set("min_price", String(params.minPrice));
  if (params.maxPrice) searchParams.set("max_price", String(params.maxPrice));

  const qs = searchParams.toString();
  const endpoint = qs ? `/products?${qs}` : "/products";
  const res = await fetchApi(endpoint);
  // Return data raw; adapt consumer to shape
  return res;
}

export async function getProduct(idOrSlug) {
  const res = await fetchApi(`/products/${idOrSlug}`);
  return res;
}

// Payments (PayOS)
export async function createPayOSCheckout(payload) {
  const response = await fetchApi(`/payment/payos/checkout`, {
    method: "POST",
    data: payload,
  });
  return response;
}
