type BadgeProps = {
  value: string;
  className?: string;
};

const badgeStyles: Record<string, string> = {
  active: 'bg-emerald-100 text-emerald-800',
  inactive: 'bg-slate-100 text-slate-700',
  entry: 'bg-emerald-100 text-emerald-800',
  exit: 'bg-rose-100 text-rose-800',
  inside: 'bg-sky-100 text-sky-800',
  outside: 'bg-orange-100 text-orange-800',
};

export function StatusBadge({ value, className = '' }: BadgeProps) {
  const normalized = value?.toLowerCase() ?? '';
  const variant = badgeStyles[normalized] ?? 'bg-slate-100 text-slate-700';

  return (
    <span className={`inline-flex items-center rounded-full px-2.5 py-1 text-[11px] font-semibold uppercase tracking-[0.18em] ${variant} ${className}`}>
      {value}
    </span>
  );
}
