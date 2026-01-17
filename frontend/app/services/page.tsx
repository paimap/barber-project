import ServicesClient from "./ServicesClient"
import { fetchWithAuth } from "@/lib/fetchWithAuth"
import { redirect } from "next/navigation"
import { ProductType, ServiceType, ProductServiceSummary } from "./types";

async function getProducts() {
    const res = await fetchWithAuth("http://localhost:8080/api/products");
    if (!res.ok) return [];
    const data = await res.json();
    return data.products as ProductType[];
}

async function getServiceTypes() {
    const res = await fetchWithAuth("http://localhost:8080/api/service-type");
    if (!res.ok) return [];
    const data = await res.json();
    return data.servicesType as ServiceType[];
}

// Fungsi fetch baru untuk Summary
async function getSummary() {
    const res = await fetchWithAuth("http://localhost:8080/api/admin/product-service/summary");
    if (!res.ok) return { product_revenue: 0, service_revenue: 0, product_sold: 0, service_performed: 0 };
    const json = await res.json();
    return json.data as ProductServiceSummary;
}

export default async function ServicesPage() {
    // Jalankan semua fetch secara paralel
    const [products, services, summary] = await Promise.all([
        getProducts(),
        getServiceTypes(),
        getSummary()
    ]);

    return (
        <ServicesClient 
            productData={products} 
            serviceData={services}
            summaryData={summary} 
        />
    );
}