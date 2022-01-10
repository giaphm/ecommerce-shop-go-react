export function formatDate(d: Date) {
    return d.getFullYear().toString() + "-" + ((d.getMonth() + 1).toString().length == 2 ? (d.getMonth() + 1).toString() : "0" + (d.getMonth() + 1).toString()) + "-" + (d.getDate().toString().length == 2 ? d.getDate().toString() : "0" + d.getDate().toString());
}

export function formatHour(d: Date) {
    return (d.getHours().toString().length == 2 ?
    d.getHours().toString() :
    "0" + d.getHours().toString()) + ":" + (d.getMinutes().toString().length == 2 ?
        d.getMinutes().toString() :
        "0" + d.getMinutes().toString());
}

export function formatDateTime(d: Date) {
    return formatDate(d) + " " + formatHour(d)
}

