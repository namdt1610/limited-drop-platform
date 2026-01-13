/**
 * Vietnam Addresses Service
 *
 * Loads Vietnamese provinces, districts, and wards from provinces.open-api.vn
 */

const VIETNAM_PROVINCES_API = "https://provinces.open-api.vn/api";

let provincesCache = [];
let districtsCache = {};
let wardsCache = {};

/**
 * Format division type (e.g., "tinh" -> "Tỉnh", "thanh pho" -> "Thành phố")
 */
function formatDivisionType(divisionType) {
  if (!divisionType) return "";
  return divisionType
    .split(" ")
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");
}

/**
 * Build label from name and division type
 */
function buildLabel(name, divisionType) {
  return name;
}

/**
 * Fetch all provinces
 */
export async function getProvinces() {
  if (provincesCache.length > 0) {
    return provincesCache;
  }

  try {
    const response = await fetch(`${VIETNAM_PROVINCES_API}/?depth=1`);
    if (!response.ok)
      throw new Error(`Failed to fetch provinces: ${response.status}`);

    const data = await response.json();
    provincesCache = data.map((province) => ({
      value: String(province.code),
      code: String(province.code),
      name: province.name,
      divisionType: formatDivisionType(province.division_type),
      label: buildLabel(province.name, province.division_type),
    }));

    return provincesCache;
  } catch (error) {
    return [];
  }
}

/**
 * Fetch districts for a given province
 */
export async function getDistricts(provinceCode) {
  if (!provinceCode) return [];

  if (districtsCache[provinceCode]) {
    return districtsCache[provinceCode];
  }

  try {
    const response = await fetch(
      `${VIETNAM_PROVINCES_API}/p/${provinceCode}?depth=2`
    );
    if (!response.ok) {
      throw new Error(`Failed to fetch districts for province ${provinceCode}`);
    }

    const data = await response.json();
    const mapped = (data.districts || []).map((district) => ({
      value: String(district.code),
      code: String(district.code),
      name: district.name,
      divisionType: formatDivisionType(district.division_type),
      label: buildLabel(district.name, district.division_type),
    }));

    districtsCache[provinceCode] = mapped;
    return mapped;
  } catch (error) {
    return [];
  }
}

/**
 * Fetch wards for a given district
 */
export async function getWards(districtCode) {
  if (!districtCode) return [];

  if (wardsCache[districtCode]) {
    return wardsCache[districtCode];
  }

  try {
    const response = await fetch(
      `${VIETNAM_PROVINCES_API}/d/${districtCode}?depth=2`
    );
    if (!response.ok) {
      throw new Error(`Failed to fetch wards for district ${districtCode}`);
    }

    const data = await response.json();
    const mapped = (data.wards || []).map((ward) => ({
      value: String(ward.code),
      code: String(ward.code),
      name: ward.name,
      divisionType: formatDivisionType(ward.division_type),
      label: buildLabel(ward.name, ward.division_type),
    }));

    wardsCache[districtCode] = mapped;
    return mapped;
  } catch (error) {
    return [];
  }
}

/**
 * Clear all caches
 */
export function clearAddressCaches() {
  provincesCache = [];
  districtsCache = {};
  wardsCache = {};
}
