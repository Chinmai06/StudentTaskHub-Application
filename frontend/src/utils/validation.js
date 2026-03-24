export function isUflEmail(email) {
  return /^[A-Za-z0-9._%+-]+@ufl\.edu$/i.test(email.trim());
}

export function isValidUfid(ufid) {
  return /^\d{8}$/.test(ufid.trim());
}

export function isNonEmptyText(value) {
  return value.trim().length > 0;
}
