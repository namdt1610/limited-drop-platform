/**
 * Vietnam Addresses - Comprehensive Table-Driven Tests
 */
import { describe, it, expect, beforeEach, vi } from 'vitest'

// Mock fetch globally
const mockFetch = vi.fn()
global.fetch = mockFetch

// Import after mocking
const {
    getProvinces,
    getDistricts,
    getWards,
    clearAddressCaches,
} = await import('./vietnam-addresses.js')

describe('Vietnam Addresses Service', () => {
    beforeEach(() => {
        mockFetch.mockReset()
        clearAddressCaches()
    })

    // =============================================================================
    // GET PROVINCES TESTS
    // =============================================================================
    describe('getProvinces', () => {
        const testCases = [
            {
                name: 'fetches and formats provinces',
                mockResponse: [
                    { code: 1, name: 'Ha Noi', division_type: 'thanh pho' },
                    { code: 79, name: 'Ho Chi Minh', division_type: 'thanh pho' },
                ],
                expectedLength: 2,
                expectedFirst: {
                    value: '1',
                    code: '1',
                    name: 'Ha Noi',
                    divisionType: 'Thanh Pho',
                    label: 'Ha Noi',
                },
            },
            {
                name: 'handles empty response',
                mockResponse: [],
                expectedLength: 0,
            },
        ]

        testCases.forEach(({ name, mockResponse, expectedLength, expectedFirst }) => {
            it(`should ${name}`, async () => {
                mockFetch.mockResolvedValueOnce({
                    ok: true,
                    json: () => Promise.resolve(mockResponse),
                })

                const result = await getProvinces()
                expect(result).toHaveLength(expectedLength)
                if (expectedFirst && result.length > 0) {
                    expect(result[0]).toEqual(expectedFirst)
                }
            })
        })

        it('should return cached provinces on second call', async () => {
            mockFetch.mockResolvedValueOnce({
                ok: true,
                json: () => Promise.resolve([{ code: 1, name: 'Ha Noi', division_type: 'tinh' }]),
            })

            await getProvinces()
            await getProvinces()

            expect(mockFetch).toHaveBeenCalledTimes(1)
        })

        it('should return empty array on fetch error', async () => {
            mockFetch.mockResolvedValueOnce({ ok: false, status: 500 })

            const result = await getProvinces()
            expect(result).toEqual([])
        })
    })

    // =============================================================================
    // GET DISTRICTS TESTS
    // =============================================================================
    describe('getDistricts', () => {
        const testCases = [
            {
                name: 'fetches districts for province',
                provinceCode: '1',
                mockResponse: {
                    districts: [
                        { code: 1, name: 'Ba Dinh', division_type: 'quan' },
                        { code: 2, name: 'Hoan Kiem', division_type: 'quan' },
                    ],
                },
                expectedLength: 2,
            },
            {
                name: 'returns empty for null province code',
                provinceCode: null,
                mockResponse: null,
                expectedLength: 0,
                skipFetch: true,
            },
            {
                name: 'returns empty for empty province code',
                provinceCode: '',
                mockResponse: null,
                expectedLength: 0,
                skipFetch: true,
            },
            {
                name: 'handles missing districts array',
                provinceCode: '99',
                mockResponse: {},
                expectedLength: 0,
            },
        ]

        testCases.forEach(({ name, provinceCode, mockResponse, expectedLength, skipFetch }) => {
            it(`should ${name}`, async () => {
                if (!skipFetch) {
                    mockFetch.mockResolvedValueOnce({
                        ok: true,
                        json: () => Promise.resolve(mockResponse),
                    })
                }

                const result = await getDistricts(provinceCode)
                expect(result).toHaveLength(expectedLength)
            })
        })

        it('should cache districts by province code', async () => {
            mockFetch.mockResolvedValueOnce({
                ok: true,
                json: () => Promise.resolve({ districts: [{ code: 1, name: 'Test', division_type: 'quan' }] }),
            })

            await getDistricts('1')
            await getDistricts('1')

            expect(mockFetch).toHaveBeenCalledTimes(1)
        })
    })

    // =============================================================================
    // GET WARDS TESTS
    // =============================================================================
    describe('getWards', () => {
        const testCases = [
            {
                name: 'fetches wards for district',
                districtCode: '1',
                mockResponse: {
                    wards: [
                        { code: 1, name: 'Phuong 1', division_type: 'phuong' },
                        { code: 2, name: 'Phuong 2', division_type: 'phuong' },
                    ],
                },
                expectedLength: 2,
            },
            {
                name: 'returns empty for null district code',
                districtCode: null,
                mockResponse: null,
                expectedLength: 0,
                skipFetch: true,
            },
            {
                name: 'returns empty for empty district code',
                districtCode: '',
                mockResponse: null,
                expectedLength: 0,
                skipFetch: true,
            },
            {
                name: 'handles missing wards array',
                districtCode: '99',
                mockResponse: {},
                expectedLength: 0,
            },
        ]

        testCases.forEach(({ name, districtCode, mockResponse, expectedLength, skipFetch }) => {
            it(`should ${name}`, async () => {
                if (!skipFetch) {
                    mockFetch.mockResolvedValueOnce({
                        ok: true,
                        json: () => Promise.resolve(mockResponse),
                    })
                }

                const result = await getWards(districtCode)
                expect(result).toHaveLength(expectedLength)
            })
        })

        it('should cache wards by district code', async () => {
            mockFetch.mockResolvedValueOnce({
                ok: true,
                json: () => Promise.resolve({ wards: [{ code: 1, name: 'Test', division_type: 'phuong' }] }),
            })

            await getWards('1')
            await getWards('1')

            expect(mockFetch).toHaveBeenCalledTimes(1)
        })
    })

    // =============================================================================
    // CLEAR CACHES TEST
    // =============================================================================
    describe('clearAddressCaches', () => {
        it('should clear all caches', async () => {
            // Populate caches
            mockFetch.mockResolvedValueOnce({
                ok: true,
                json: () => Promise.resolve([{ code: 1, name: 'Test', division_type: 'tinh' }]),
            })
            await getProvinces()

            // Clear caches
            clearAddressCaches()

            // Should fetch again
            mockFetch.mockResolvedValueOnce({
                ok: true,
                json: () => Promise.resolve([{ code: 2, name: 'Test2', division_type: 'tinh' }]),
            })

            const result = await getProvinces()
            expect(result[0].code).toBe('2')
            expect(mockFetch).toHaveBeenCalledTimes(2)
        })
    })
})
