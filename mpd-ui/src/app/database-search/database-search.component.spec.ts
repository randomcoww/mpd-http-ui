import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DatabaseSearchComponent } from './database-search.component';

describe('DatabaseSearchComponent', () => {
  let component: DatabaseSearchComponent;
  let fixture: ComponentFixture<DatabaseSearchComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ DatabaseSearchComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DatabaseSearchComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
