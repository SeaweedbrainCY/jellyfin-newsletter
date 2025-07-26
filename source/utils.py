import datetime as dt
from configuration import logging


def summarize_ranges(nums):
    """
    Summarizes a list of integers into ranges.
    For example, [0, 1, 2, 4, 5, 7] becomes ['0-2', '4-5', '7'].
    """
    if not nums:
        return []
    # convert all elements to integers
    try:
        nums = list(map(int, nums))
    except Exception as e:
        logging.error(f"Error while checking episodes for a show. Episodes list will not be displayed in the final email due to this error : {e}")
        return None
    nums = sorted(nums)
    result = []
    start = nums[0]
    end = nums[0]

    for n in nums[1:]:
        if n == end + 1:
            end = n
        else:
            if start == end:
                result.append(str(start))
            else:
                result.append(f"{start}-{end}")
            start = end = n

    if start == end:
        result.append(str(start))
    else:
        result.append(f"{start}-{end}")

    return result


def get_last_newsletter_date():
    """
    Returns the date of the last newsletter.
    If the file does not exist, it returns None.
    """
    try:
        with open("./config/LAST_NEWSLETTER.txt", "r") as f:
            date_str = f.read().strip()
        try:
            return dt.datetime.fromisoformat(date_str)
        except ValueError:
            logging.error(f"Error while parsing the date from LAST_NEWSLETTER.txt. Expected ISO format, got: {date_str}. It is highly recommended to delete this file and let the program create a new one.")
            return None
    except FileNotFoundError:
        return None

def save_last_newsletter_date(date):
    """
    Saves the date of the last newsletter to a file.
    The date should be a datetime object.
    """
    with open("./config/LAST_NEWSLETTER.txt", "w") as f:
        f.write(date.isoformat())