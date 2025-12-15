export default function formatFullDate(date) {
    const d = new Date(date);

    const day = d.getDate();
    const month = d.getMonth() + 1;
    const year = d.getFullYear();

    const hours = d.getHours().toString().padStart(2, '0');
    const minutes = d.getMinutes().toString().padStart(2, '0');

    return `${year}/${month}/${day} - ${hours}:${minutes}`;
}
