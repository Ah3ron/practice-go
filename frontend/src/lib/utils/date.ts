import { formatDistanceToNow } from 'date-fns';

/**
 * Safely formats a date using formatDistanceToNow with proper error handling
 * @param dateValue - The date value to format (string, Date, null, or undefined)
 * @returns Formatted date string or fallback text for invalid dates
 */
export function safeFormatDistanceToNow(dateValue: string | Date | null | undefined): string {
  if (!dateValue) return 'Unknown';

  try {
    const date = new Date(dateValue);
    // Check if date is valid
    if (isNaN(date.getTime())) {
      return 'Invalid date';
    }
    return formatDistanceToNow(date, { addSuffix: true });
  } catch (error) {
    return 'Invalid date';
  }
}

/**
 * Safely formats a date to a local date string
 * @param dateValue - The date value to format
 * @returns Formatted date string or fallback text for invalid dates
 */
export function safeFormatDate(dateValue: string | Date | null | undefined): string {
  if (!dateValue) return 'Unknown';

  try {
    const date = new Date(dateValue);
    if (isNaN(date.getTime())) {
      return 'Invalid date';
    }
    return date.toLocaleDateString();
  } catch (error) {
    return 'Invalid date';
  }
}

/**
 * Checks if a date value is valid
 * @param dateValue - The date value to validate
 * @returns True if the date is valid, false otherwise
 */
export function isValidDate(dateValue: string | Date | null | undefined): boolean {
  if (!dateValue) return false;

  try {
    const date = new Date(dateValue);
    return !isNaN(date.getTime());
  } catch (error) {
    return false;
  }
}
