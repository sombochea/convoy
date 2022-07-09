import { CommonModule } from '@angular/common';
import { Component, Input, OnInit } from '@angular/core';

@Component({
	selector: 'convoy-card',
	standalone: true,
	imports: [CommonModule],
	templateUrl: './card.component.html',
	styleUrls: ['./card.component.scss']
})
export class CardComponent implements OnInit {
	@Input('class') class!: string;
  
	constructor() {}

	ngOnInit(): void {}
}
