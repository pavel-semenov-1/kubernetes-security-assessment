import { Item } from "./Item";
import { Misconfiguration } from "./Misconfiguration";

export type ItemListProps<T extends Item> = {
    data: Record<string, T[]> | null;
    onClickResolve: ((category: string, index: number, item: Misconfiguration) => void) | undefined;
    onClickDelete: (category: string, index: number, item: T) => void;
    openCategory: number | null;
    setOpenCategory: (category: number | null) => void;
    openItem: number | null;
    setOpenItem: (category: number | null) => void;
}