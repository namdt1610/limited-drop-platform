/**
 * API Utils - Comprehensive Table-Driven Tests
 * Uses mocked fetch for all API calls
 */
import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'

// Mock fetch globally
const mockFetch = vi.fn()
global.fetch = mockFetch

// Mock import.meta.env
vi.stubGlobal('import', { meta: { env: { VITE_API_URL: 'http://test-api.com' } } })

// Re-import after mocking
const apiModule = await import('./api.js')
const {
    fetchApi,
    getDrops,
    getDropStatus,
    purchaseDrop,
    trackOrder,
    verifyPayment,
    cancelPayment,
    verifySymbicode,
    getProducts,
    getProduct,
} = apiModule

describe('fetchApi', () => {
    beforeEach(() => {
        mockFetch.mockReset()
    })

    const successCases = [
        {
            name: 'GET request returns JSON',
            endpoint: '/test',
            options: {},
            mockResponse: { data: 'test' },
            expected: { data: 'test' },
        },
        {
            name: 'POST request with data',
            endpoint: '/test',
            options: { method: 'POST', data: { foo: 'bar' } },
            mockResponse: { success: true },
            expected: { success: true },
        },
        {
            name: 'empty response returns empty object',
            endpoint: '/empty',
            options: {},
            mockResponse: null, // Will throw on json parse
            expected: {},
        },
    ]

    successCases.forEach(({ name, endpoint, options, mockResponse, expected }) => {
        it(`should handle ${name}`, async () => {
            mockFetch.mockResolvedValueOnce({
                ok: true,
                json: mockResponse ? () => Promise.resolve(mockResponse) : () => Promise.reject(),
            })

            const result = await fetchApi(endpoint, options)
            expect(result).toEqual(expected)
        })
    })

    const errorCases = [
        {
            name: 'API error with JSON message',
            status: 400,
            errorBody: { message: 'Bad Request' },
            expectedError: 'Bad Request',
        },
        {
            name: 'API error with error field',
            status: 404,
            errorBody: { error: 'Not Found' },
            expectedError: 'Not Found',
        },
        {
            name: 'API error with text response',
            status: 500,
            errorBody: null,
            textBody: 'Internal Server Error',
            expectedError: 'Internal Server Error',
        },
        {
            name: 'API error with status only',
            status: 403,
            errorBody: {},
            expectedError: 'API Error: 403',
        },
    ]

    errorCases.forEach(({ name, status, errorBody, textBody, expectedError }) => {
        it(`should handle ${name}`, async () => {
            mockFetch.mockResolvedValueOnce({
                ok: false,
                status,
                json: errorBody ? () => Promise.resolve(errorBody) : () => Promise.reject(),
                text: () => Promise.resolve(textBody || ''),
            })

            await expect(fetchApi('/error')).rejects.toThrow(expectedError)
        })
    })
})

describe('getDrops', () => {
    beforeEach(() => mockFetch.mockReset())

    const testCases = [
        {
            name: 'returns array directly',
            response: [{ id: 1 }, { id: 2 }],
            expected: [{ id: 1 }, { id: 2 }],
        },
        {
            name: 'returns drops from object',
            response: { drops: [{ id: 3 }] },
            expected: [{ id: 3 }],
        },
        {
            name: 'returns empty array for non-array',
            response: { other: 'data' },
            expected: [],
        },
    ]

    testCases.forEach(({ name, response, expected }) => {
        it(`should handle ${name}`, async () => {
            mockFetch.mockResolvedValueOnce({
                ok: true,
                json: () => Promise.resolve(response),
            })

            const result = await getDrops()
            expect(result).toEqual(expected)
        })
    })
})

describe('getDropStatus', () => {
    beforeEach(() => mockFetch.mockReset())

    it('should fetch drop status by ID', async () => {
        const mockStatus = { drop_id: 1, available: 50, phase: 'LIVE' }
        mockFetch.mockResolvedValueOnce({
            ok: true,
            json: () => Promise.resolve(mockStatus),
        })

        const result = await getDropStatus(1)
        expect(result).toEqual(mockStatus)
        expect(mockFetch).toHaveBeenCalledWith(
            expect.stringContaining('/api/drops/1/status'),
            expect.any(Object)
        )
    })
})

describe('purchaseDrop', () => {
    beforeEach(() => mockFetch.mockReset())

    const testCases = [
        {
            name: 'successful purchase',
            dropId: 1,
            data: { phone: '0123456789', email: 'test@test.com' },
            response: { payment_url: 'https://pay.os/xxx', order_code: 123 },
        },
        {
            name: 'purchase with full data',
            dropId: 2,
            data: {
                quantity: 1,
                name: 'John',
                phone: '0987654321',
                email: 'john@test.com',
                address: '123 St',
                province: 'HCM',
                district: 'D1',
                ward: 'W1',
            },
            response: { payment_url: 'https://pay.os/yyy', order_code: 456 },
        },
    ]

    testCases.forEach(({ name, dropId, data, response }) => {
        it(`should handle ${name}`, async () => {
            mockFetch.mockResolvedValueOnce({
                ok: true,
                json: () => Promise.resolve(response),
            })

            const result = await purchaseDrop(dropId, data)
            expect(result).toEqual(response)
            expect(mockFetch).toHaveBeenCalledWith(
                expect.stringContaining(`/api/drops/${dropId}/purchase`),
                expect.objectContaining({ method: 'POST' })
            )
        })
    })
})

describe('verifySymbicode', () => {
    beforeEach(() => mockFetch.mockReset())

    const testCases = [
        {
            name: 'first activation',
            code: '550e8400-e29b-41d4-a716-446655440000',
            response: { symbicode: { id: 1 }, is_first_activation: true },
        },
        {
            name: 'already activated',
            code: '550e8400-e29b-41d4-a716-446655440001',
            response: { symbicode: { id: 2 }, is_first_activation: false },
        },
    ]

    testCases.forEach(({ name, code, response }) => {
        it(`should handle ${name}`, async () => {
            mockFetch.mockResolvedValueOnce({
                ok: true,
                json: () => Promise.resolve(response),
            })

            const result = await verifySymbicode(code)
            expect(result).toEqual(response)
        })
    })
})

describe('getProducts', () => {
    beforeEach(() => mockFetch.mockReset())

    const testCases = [
        {
            name: 'no params',
            params: {},
            expectedEndpoint: '/products',
        },
        {
            name: 'with page',
            params: { page: 2 },
            expectedEndpoint: '/products?page=2',
        },
        {
            name: 'with all params',
            params: { page: 1, limit: 10, search: 'test', sort: 'price', minPrice: 100, maxPrice: 500 },
            expectedEndpoint: '/products?page=1&limit=10&search=test&sort=price&min_price=100&max_price=500',
        },
    ]

    testCases.forEach(({ name, params, expectedEndpoint }) => {
        it(`should build URL for ${name}`, async () => {
            mockFetch.mockResolvedValueOnce({
                ok: true,
                json: () => Promise.resolve([]),
            })

            await getProducts(params)
            expect(mockFetch).toHaveBeenCalledWith(
                expect.stringContaining(expectedEndpoint),
                expect.any(Object)
            )
        })
    })
})

describe('getProduct', () => {
    beforeEach(() => mockFetch.mockReset())

    const testCases = [
        { name: 'by ID', idOrSlug: 1 },
        { name: 'by slug', idOrSlug: 'test-product' },
    ]

    testCases.forEach(({ name, idOrSlug }) => {
        it(`should fetch product ${name}`, async () => {
            mockFetch.mockResolvedValueOnce({
                ok: true,
                json: () => Promise.resolve({ id: 1, name: 'Test' }),
            })

            const result = await getProduct(idOrSlug)
            expect(result).toHaveProperty('id')
            expect(mockFetch).toHaveBeenCalledWith(
                expect.stringContaining(`/products/${idOrSlug}`),
                expect.any(Object)
            )
        })
    })
})
