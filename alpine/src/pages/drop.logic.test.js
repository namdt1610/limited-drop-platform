/**
 * Drop Logic - Comprehensive Table-Driven Tests
 * Tests the core drop purchase flow logic
 */
import { describe, it, expect, beforeEach, vi } from 'vitest'

// =============================================================================
// VALIDATION LOGIC TESTS (extracted from DropLogic)
// =============================================================================

// Re-implement validation logic for testing (pure functions)
function validatePhone(phone) {
    if (!phone) return 'Truong nay la bat buoc'
    if (!/^(0|\+84)[3|5|7|8|9][0-9]{8}$/.test(phone)) {
        return 'So dien thoai khong hop le'
    }
    return null
}

function validateEmail(email) {
    if (!email) return 'Truong nay la bat buoc'
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
        return 'Email khong hop le'
    }
    return null
}

function validateRequired(value) {
    if (!value) return 'Truong nay la bat buoc'
    return null
}

function maskPhone(phone) {
    if (!phone) return ''
    return phone.replace(/(\d{3})\d{3}(\d{3})/, '$1***$2')
}

function determinePhase(dropData, serverOffset = 0) {
    if (!dropData) return 'WAITING'

    const now = new Date(Date.now() + serverOffset)
    const start = new Date(dropData.starts_at)
    const end = dropData.ends_at ? new Date(dropData.ends_at) : null
    const isSoldOut = Number(dropData.available) <= 0

    if (now < start) return 'WAITING'
    if (isSoldOut) return 'SOLD_OUT'
    if (end && now >= end) return 'ENDED'
    return 'LIVE'
}

function isDisabled(phase, isSoldOut) {
    return !(phase === 'LIVE' && !isSoldOut)
}

describe('Phone Validation', () => {
    const testCases = [
        // Valid phones
        { name: 'valid 09x', input: '0912345678', expected: null },
        { name: 'valid 03x', input: '0358765432', expected: null },
        { name: 'valid 05x', input: '0567890123', expected: null },
        { name: 'valid 07x', input: '0789012345', expected: null },
        { name: 'valid 08x', input: '0891234567', expected: null },
        { name: 'valid +84', input: '+84912345678', expected: null },

        // Invalid phones
        { name: 'empty', input: '', expected: 'Truong nay la bat buoc' },
        { name: 'too short', input: '091234567', expected: 'So dien thoai khong hop le' },
        { name: 'too long', input: '09123456789', expected: 'So dien thoai khong hop le' },
        { name: 'invalid prefix', input: '0112345678', expected: 'So dien thoai khong hop le' },
        { name: 'with spaces', input: '091 234 5678', expected: 'So dien thoai khong hop le' },
        { name: 'with dashes', input: '091-234-5678', expected: 'So dien thoai khong hop le' },
        { name: 'letters', input: 'abc1234567', expected: 'So dien thoai khong hop le' },
    ]

    testCases.forEach(({ name, input, expected }) => {
        it(`should validate ${name}`, () => {
            expect(validatePhone(input)).toBe(expected)
        })
    })
})

describe('Email Validation', () => {
    const testCases = [
        // Valid emails
        { name: 'simple email', input: 'test@example.com', expected: null },
        { name: 'with subdomain', input: 'user@mail.example.com', expected: null },
        { name: 'with numbers', input: 'user123@test.com', expected: null },
        { name: 'with dots', input: 'first.last@example.com', expected: null },
        { name: 'with plus', input: 'user+tag@example.com', expected: null },

        // Invalid emails
        { name: 'empty', input: '', expected: 'Truong nay la bat buoc' },
        { name: 'no @', input: 'testexample.com', expected: 'Email khong hop le' },
        { name: 'no domain', input: 'test@', expected: 'Email khong hop le' },
        { name: 'no TLD', input: 'test@example', expected: 'Email khong hop le' },
        { name: 'spaces', input: 'test @example.com', expected: 'Email khong hop le' },
        { name: 'double @', input: 'test@@example.com', expected: 'Email khong hop le' },
    ]

    testCases.forEach(({ name, input, expected }) => {
        it(`should validate ${name}`, () => {
            expect(validateEmail(input)).toBe(expected)
        })
    })
})

describe('Required Field Validation', () => {
    const testCases = [
        { name: 'non-empty string', input: 'value', expected: null },
        { name: 'empty string', input: '', expected: 'Truong nay la bat buoc' },
        { name: 'null', input: null, expected: 'Truong nay la bat buoc' },
        { name: 'undefined', input: undefined, expected: 'Truong nay la bat buoc' },
        { name: 'whitespace only', input: '   ', expected: null }, // Not trimmed in original
        { name: 'zero', input: 0, expected: 'Truong nay la bat buoc' },
    ]

    testCases.forEach(({ name, input, expected }) => {
        it(`should validate ${name}`, () => {
            expect(validateRequired(input)).toBe(expected)
        })
    })
})

describe('Phone Masking', () => {
    const testCases = [
        { name: 'standard 10-digit', input: '0912345678', expected: '091***5678' },
        { name: '+84 format', input: '+84912345678', expected: '+849***45678' },
        { name: 'empty', input: '', expected: '' },
        { name: 'null', input: null, expected: '' },
        { name: 'short number', input: '12345', expected: '12345' }, // No mask
    ]

    testCases.forEach(({ name, input, expected }) => {
        it(`should mask ${name}`, () => {
            expect(maskPhone(input)).toBe(expected)
        })
    })
})

describe('Phase Determination', () => {
    const now = Date.now()

    const testCases = [
        {
            name: 'WAITING - before start',
            dropData: {
                starts_at: new Date(now + 3600000).toISOString(), // 1 hour later
                ends_at: null,
                available: 100,
            },
            expected: 'WAITING',
        },
        {
            name: 'LIVE - after start, before end, stock available',
            dropData: {
                starts_at: new Date(now - 3600000).toISOString(), // 1 hour ago
                ends_at: new Date(now + 3600000).toISOString(), // 1 hour later
                available: 50,
            },
            expected: 'LIVE',
        },
        {
            name: 'LIVE - no end time, stock available',
            dropData: {
                starts_at: new Date(now - 3600000).toISOString(),
                ends_at: null,
                available: 50,
            },
            expected: 'LIVE',
        },
        {
            name: 'SOLD_OUT - no stock',
            dropData: {
                starts_at: new Date(now - 3600000).toISOString(),
                ends_at: new Date(now + 3600000).toISOString(),
                available: 0,
            },
            expected: 'SOLD_OUT',
        },
        {
            name: 'ENDED - past end time',
            dropData: {
                starts_at: new Date(now - 7200000).toISOString(), // 2 hours ago
                ends_at: new Date(now - 3600000).toISOString(), // 1 hour ago
                available: 50,
            },
            expected: 'ENDED',
        },
        {
            name: 'null drop data',
            dropData: null,
            expected: 'WAITING',
        },
    ]

    testCases.forEach(({ name, dropData, expected }) => {
        it(`should determine ${name}`, () => {
            expect(determinePhase(dropData)).toBe(expected)
        })
    })
})

describe('Button Disabled State', () => {
    const testCases = [
        { name: 'LIVE + not sold out', phase: 'LIVE', isSoldOut: false, expected: false },
        { name: 'LIVE + sold out', phase: 'LIVE', isSoldOut: true, expected: true },
        { name: 'WAITING', phase: 'WAITING', isSoldOut: false, expected: true },
        { name: 'ENDED', phase: 'ENDED', isSoldOut: false, expected: true },
        { name: 'SOLD_OUT phase', phase: 'SOLD_OUT', isSoldOut: true, expected: true },
    ]

    testCases.forEach(({ name, phase, isSoldOut, expected }) => {
        it(`should be disabled=${expected} for ${name}`, () => {
            expect(isDisabled(phase, isSoldOut)).toBe(expected)
        })
    })
})

// =============================================================================
// FOMO STATUS TESTS
// =============================================================================

function updateFomoStatuses(contact, phase, isSoldOut, winner, fomoTick) {
    const base = [
        { phone: '09x xxx 888', action: 'dang thanh toan', time: '0.4s' },
        { phone: '08x xxx 555', action: 'vua giu slot', time: '0.7s' },
        { phone: '03x xxx 123', action: 'dang nhap dia chi', time: '0.9s' },
        { phone: '07x xxx 246', action: 'dang kiem tra thong tin', time: '1.1s' },
    ]

    const user =
        contact.phone && phase === 'LIVE'
            ? [{ phone: maskPhone(contact.phone), action: 'dang tranh slot', time: '...' }]
            : []

    const winnerStatus =
        isSoldOut && winner
            ? [{ phone: winner.maskedPhone, action: 'da chot', time: winner.time || '-' }]
            : []

    const all = [...winnerStatus, ...base, ...user]
    if (!all.length) return []

    const start = fomoTick % all.length
    return [...all.slice(start), ...all.slice(0, start)].slice(0, Math.min(4, all.length))
}

describe('FOMO Status Generation', () => {
    const testCases = [
        {
            name: 'base statuses only',
            contact: { phone: '' },
            phase: 'WAITING',
            isSoldOut: false,
            winner: null,
            fomoTick: 0,
            expectedLength: 4,
        },
        {
            name: 'includes user when LIVE with phone',
            contact: { phone: '0912345678' },
            phase: 'LIVE',
            isSoldOut: false,
            winner: null,
            fomoTick: 0,
            expectedLength: 4, // max 4
        },
        {
            name: 'includes winner when sold out',
            contact: { phone: '' },
            phase: 'SOLD_OUT',
            isSoldOut: true,
            winner: { maskedPhone: '091***678', time: '2.5s' },
            fomoTick: 0,
            expectedLength: 4,
        },
        {
            name: 'rotates based on fomoTick',
            contact: { phone: '' },
            phase: 'WAITING',
            isSoldOut: false,
            winner: null,
            fomoTick: 1,
            expectedLength: 4,
        },
    ]

    testCases.forEach(({ name, contact, phase, isSoldOut, winner, fomoTick, expectedLength }) => {
        it(`should generate ${name}`, () => {
            const result = updateFomoStatuses(contact, phase, isSoldOut, winner, fomoTick)
            expect(result).toHaveLength(expectedLength)
        })
    })
})

// =============================================================================
// COMPLETE VALIDATION TESTS
// =============================================================================

function validateAll(contact) {
    const errors = {}
    const fields = ['phone', 'email', 'name', 'address', 'province', 'district', 'ward']

    fields.forEach((field) => {
        const value = contact[field]
        if (!value) {
            errors[field] = 'Truong nay la bat buoc'
        } else if (field === 'phone') {
            const phoneErr = validatePhone(value)
            if (phoneErr) errors[field] = phoneErr
        } else if (field === 'email') {
            const emailErr = validateEmail(value)
            if (emailErr) errors[field] = emailErr
        }
    })

    return { isValid: Object.keys(errors).length === 0, errors }
}

describe('Complete Form Validation', () => {
    const testCases = [
        {
            name: 'valid complete form',
            contact: {
                phone: '0912345678',
                email: 'test@example.com',
                name: 'John Doe',
                address: '123 Main St',
                province: 'Ho Chi Minh',
                district: 'District 1',
                ward: 'Ward 1',
            },
            expectedValid: true,
            expectedErrorCount: 0,
        },
        {
            name: 'empty form',
            contact: {
                phone: '',
                email: '',
                name: '',
                address: '',
                province: '',
                district: '',
                ward: '',
            },
            expectedValid: false,
            expectedErrorCount: 7,
        },
        {
            name: 'invalid phone only',
            contact: {
                phone: '123',
                email: 'test@example.com',
                name: 'John',
                address: '123 St',
                province: 'HCM',
                district: 'D1',
                ward: 'W1',
            },
            expectedValid: false,
            expectedErrorCount: 1,
        },
        {
            name: 'invalid email only',
            contact: {
                phone: '0912345678',
                email: 'invalid',
                name: 'John',
                address: '123 St',
                province: 'HCM',
                district: 'D1',
                ward: 'W1',
            },
            expectedValid: false,
            expectedErrorCount: 1,
        },
        {
            name: 'missing address fields',
            contact: {
                phone: '0912345678',
                email: 'test@example.com',
                name: 'John',
                address: '123 St',
                province: '',
                district: '',
                ward: '',
            },
            expectedValid: false,
            expectedErrorCount: 3,
        },
    ]

    testCases.forEach(({ name, contact, expectedValid, expectedErrorCount }) => {
        it(`should validate ${name}`, () => {
            const { isValid, errors } = validateAll(contact)
            expect(isValid).toBe(expectedValid)
            expect(Object.keys(errors)).toHaveLength(expectedErrorCount)
        })
    })
})
