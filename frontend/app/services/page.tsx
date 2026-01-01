import ServicesClient from "./ServicesClient"
import { fetchWithAuth } from "@/lib/fetchWithAuth"
import { redirect } from "next/navigation"
import { ProductType, ServiceType } from "./types";

// Fungsi helper dipisah agar clean
async function getProducts() {
    const res = await fetchWithAuth("http://localhost:8080/api/products");
    if (!res.ok) return redirect("/login");
    const data = await res.json();
    return data.products as ProductType[];
}

async function getServiceTypes() {
    const res = await fetchWithAuth("http://localhost:8080/api/service-type");
    if (!res.ok) return redirect("/login");
    const data = await res.json();
    return data.servicesType as ServiceType[];
}

export default async function MitraPage() {
    const [products, services] = await Promise.all([
        getProducts(),
        getServiceTypes()
    ]);

    return (
        <ServicesClient 
            productData={products} 
            serviceData={services} // Kirim data service jika ada
        />
    );
}