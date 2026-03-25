import { describe, expect, it } from 'vitest';
import { isNonEmptyText, isUflEmail, isValidUfid } from './validation';

describe('Issue #22 helpers - validation', () => {
  it('accepts a valid @ufl.edu email', () => {
    expect(isUflEmail('student@ufl.edu')).toBe(true);
  });

  it('rejects a non-UFL email', () => {
    expect(isUflEmail('student@gmail.com')).toBe(false);
  });

  it('accepts a valid 8 digit UFID', () => {
    expect(isValidUfid('12345678')).toBe(true);
  });

  it('rejects an invalid UFID', () => {
    expect(isValidUfid('1234')).toBe(false);
  });

  it('accepts non-empty trimmed text', () => {
    expect(isNonEmptyText(' Library ')).toBe(true);
  });

  it('rejects empty trimmed text', () => {
    expect(isNonEmptyText('   ')).toBe(false);
  });
});