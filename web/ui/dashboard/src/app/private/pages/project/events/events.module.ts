import { NgModule } from '@angular/core';
import { CommonModule, DatePipe } from '@angular/common';
import { EventsComponent } from './events.component';
import { Routes, RouterModule } from '@angular/router';
import { ReactiveFormsModule } from '@angular/forms';
import { DateFilterModule } from 'src/app/private/components/date-filter/date-filter.module';
import { LoaderModule } from 'src/app/private/components/loader/loader.module';
import { EventModule } from './event/event.module';
import { EventDeliveriesModule } from './event-deliveries/event-deliveries.module';
import { DropdownComponent } from 'src/app/components/dropdown/dropdown.component';
import { ButtonComponent } from 'src/app/components/button/button.component';
import { ListItemComponent } from 'src/app/components/list-item/list-item.component';
import { CardComponent } from 'src/app/components/card/card.component';
import { ChartComponent } from 'src/app/components/chart/chart.component';
import { SkeletonLoaderComponent } from 'src/app/components/skeleton-loader/skeleton-loader.component';

const routes: Routes = [{ path: '', component: EventsComponent }];

@NgModule({
	declarations: [EventsComponent],
	imports: [
		CommonModule,
		ReactiveFormsModule,
		DateFilterModule,
		RouterModule.forChild(routes),
		EventModule,
		EventDeliveriesModule,
		DropdownComponent,
		ButtonComponent,
		ListItemComponent,
		CardComponent,
        ChartComponent,
        SkeletonLoaderComponent
	],
	providers: [DatePipe]
})
export class EventsModule {}
